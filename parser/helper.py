from bs4 import BeautifulSoup
import pymongo

import requests

from settings import *
from anti_block import *

import json
import datetime
import time
import random
import logging

# configure logger
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

YANDEX_URL = 'https://realty.yandex.ru/moskva/snyat/kvartira/studiya,1,2,3,4-i-bolee-komnatnie/?priceMax=200000&sort=DATE_DESC'


class MongoDB:
    def __init__(self):
        try:
            self.db = pymongo.MongoClient(f'mongodb://{MONGO_USER}:{MONGO_PASS}@{MONGO_HOST}/{MONGO_DB}?connectTimeoutMS={MONGO_TIMEOUT}')[MONGO_DB]
            self.db.command('ismaster')
        except pymongo.errors.ServerSelectionTimeoutError:
            print("Something went wrong with MongoDB.")
            raise

    def exists(self, data):
        return self.db.get_collection(MONGO_FLATS).count_documents(data) != 0

    def collection(self, db_name):
        return self.db.get_collection(db_name)


class Parser:
    def __init__(self, db):
        self.db = db

    def update_or_do_nothing(self, data):
        document = list(self.db.collection(MONGO_FLATS).find({'url': data['url']}))[0]

        updated_fields = {}
        for key, value in document.items():
            if key in ['updated_at', 'created_at', 'updated_fields', '_id', 'processed']:
                continue

            if value != data[key]:
                updated_fields[key] = data[key]

        if updated_fields:
            if 'updated_fields' not in document:
                document['updated_fields'] = []
            document['updated_fields'].append(list(updated_fields))
            updated_fields['updated_fields'] = document['updated_fields']
            updated_fields['updated_at'] = data['updated_at']

            self.db.collection(MONGO_FLATS).update_one(
                {'url': data['url']},
                {'$set': updated_fields},
                upsert=False,
            )

    def parse_yandex(self, soup):
        """
        returns list of flats metadata
        :param soup:
        :param db:
        :return:
        """
        flats_block = str(soup.find('div', {'class', 'i-react-state i-bem'})['data-bem'])
        json_block = json.loads(flats_block)

        flats = json_block['i-react-state']['state']['search']['offers']['entities']

        for flat in flats:

            current_data = {
                'url': f'http:{flat.get("unsignedInternalUrl", "shrug")}',
                'rooms': flat.get('roomsTotal', 0),
                'address': flat.get('location', {}).get('address', ""),
                'latitude': flat.get('location', {}).get('point', {}).get('latitude', 0),
                'longitude': flat.get('location', {}).get('point', {}).get('longitude', 0),
                'price': flat.get('price', {}).get('value', 0),
                'area': flat.get('area', {}).get('value', 0),
                'fee': int(flat.get('agentFee', 0) / 100 * flat.get('price', {}).get('value', 0)),
                'prepayment': int(flat.get('prepayment', 0) / 100 * flat.get('price', {}).get('value', 0)),
                'images': flat.get('fullImages', []),
                'planImages': flat.get('extImages', {}).get('IMAGE_PLAN', {}).get('fullImages', []),
                'created_at': datetime.datetime.fromtimestamp(time.time()).ctime(),
                'updated_at': datetime.datetime.fromtimestamp(time.time()).ctime(),
                'updated_fields': [],
            }

            if current_data['url'] == 'shrug':
                continue

            for field in ('images', 'planImages'):
                current_data[field] = [url[2:] for url in current_data[field]]

            for image in current_data['planImages']:
                current_data['images'].remove(image)

            if self.db.exists({'url': current_data['url']}):
                self.update_or_do_nothing(current_data)
            else:
                self.db.collection(MONGO_FLATS).insert_one(current_data)


def configure_settings(db):
    global COOKIES
    COOKIES = [item['cookie'] for item in list(db.collection('yandex_cookies').find({}))]

    global USER_AGENTS
    USER_AGENTS = [item['user_agent'] for item in list(db.collection('user_agents').find({}))]

    global REFERERS
    REFERERS = [item['referer'] for item in list(db.collection('referers').find({}))]

    global PROXIES
    PROXIES = [item['proxy'] for item in list(db.collection('proxies').find({}))]

    global SESSIONS
    SESSIONS = [(requests.Session(), cookie) for cookie in COOKIES]


def main():
    # connect to db
    logger.info('Connecting to MongoDB...')
    parser = Parser(MongoDB())

    configure_settings(parser.db)

    counter = 1
    proxy_id = 0
    while True:
        logger.info('- Making a request...')

        # get random headers & proxy for request
        session, cookie = random.choice(SESSIONS)
        user_agent = random.choice(USER_AGENTS)
        proxy = PROXIES[proxy_id] if PROXIES else None
        proxy_id = ((proxy_id + 1) % len(PROXIES)) if PROXIES else None
        referer = random.choice(REFERERS)

        # add headers
        HEADERS['Cookie'] = cookie
        HEADERS['User-Agent'] = user_agent
        HEADERS['Referer'] = referer

        logger.info("--- Request info:")
        logger.info(f"------ User-Agent: {user_agent}")
        logger.info(f"------ Referer: {referer}")
        logger.info(f"------ Proxy: {proxy}")

        try:
            # create request
            state = session.get(YANDEX_URL, headers=HEADERS, proxies={'https': proxy} if proxy else None, timeout=7)

            # create soup for the caught state
            soup = BeautifulSoup(state.text.encode(state.encoding).decode('utf-8'), 'html.parser')

            # parse soup for the flats meta
            logger.info('- Extracting flats from html...')
            parser.parse_yandex(soup)

            logger.info(f'--------------------- Successfully completed {counter} requests ------------------')
            counter += 1
            time.sleep(28 + random.randint(3, 36))
        except Exception as e:
            logger.info('I\'m fucked up on "{}"'.format(str(e)))

            # remove cookie if we got banned
            if 'NoneType' in str(e) and 'subscriptable' in str(e):
                COOKIES.remove(cookie)
                logger.info(f'> Got rid of cookie. RIP.')

            time.sleep(17 + random.randint(12, 29))


if __name__ == "__main__":
    main()
