import json
import logging
import time

import requests

import pymongo
from bs4 import BeautifulSoup
from settings import *

YANDEX_URL = 'https://realty.yandex.ru/moskva/snyat/kvartira/studiya,1,2,3,4-i-bolee-komnatnie/?priceMax=200000&sort=DATE_DESC'

HEADERS = {
    'User-Agent': 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 YaBrowser/18.11.1.716 Yowser/2.5 Safari/537.36',
}


def exists(db, data):
    return len(list(db.get_collection(MONGO_DB).find(data))) == 0


def parse_yandex(soup, db):
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
        }

        if current_data['url'] == 'shrug':
            continue

        for field in ('images', 'planImages'):
            current_data[field] = [url[2:] for url in current_data[field]]

        for image in current_data['planImages']:
            current_data['images'].remove(image)

        if exists(db, current_data):
            db.get_collection(MONGO_DB).insert_one(current_data)


def main():
    logging.basicConfig(level=logging.INFO)
    logger = logging.getLogger(__name__)

    logger.info('Connecting to MongoDB...')
    db = pymongo.MongoClient(f'mongodb://{MONGO_USER}:{MONGO_PASS}@{MONGO_HOST}/{MONGO_DB}?connectTimeoutMS={MONGO_TIMEOUT}')[MONGO_DB]

    logger.info('Making a request...')
    state = requests.get(YANDEX_URL, headers=HEADERS)
    soup = BeautifulSoup(state.text.encode(state.encoding).decode('utf-8'), 'html.parser')

    logger.info('Extracting flats from html...')
    parse_yandex(soup, db)

    logger.info('Successfully completed.')


if __name__ == "__main__":
    while True:
        main()
        time.sleep(20)
