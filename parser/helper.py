from bs4 import BeautifulSoup

import requests
import json


YANDEX_URL = 'https://realty.yandex.ru/moskva/snyat/kvartira/studiya,1-komnatnie/?priceMax=70000&metroTransport=ON_FOOT&timeToMetro=15&priceMin=70000&includeTag=1794389&sort=DATE_DESC'

HEADERS = {
    'User-Agent': 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 YaBrowser/18.11.1.716 Yowser/2.5 Safari/537.36',
    'X-Ya-Front-Host': 'yandex.ru',
    'Connection': 'keep-alive',
    'Cache-Control': 'max-age=0',
    'Upgrade-Insecure-Requests': '1',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
    'Accept-Encoding': 'gzip, deflate, br',
    'Accept-Language': 'ru,en;q=0.9',
}


def parse_yandex(soup):
    """
        returns list of flats metadata
    :param soup:
    :return:
    """
    flats_block = str(soup.find('div', {'class', 'i-react-state i-bem'})['data-bem'])
    json_block = json.loads(flats_block)

    flats = json_block['i-react-state']['state']['search']['offers']['entities']

    print(json.dumps(flats[1], indent=2))

    data = []
    for flat in flats:

        current_data = {
            'url': f'http:{flat.get("unsignedInternalUrl", "shrug")}',
            'rooms': flat.get('roomsTotal', 0),
            'address': flat.get('location', {}).get('address', ""),
            'lat': flat.get('location', {}).get('point', {}).get('latitude', 0),
            'long': flat.get('location', {}).get('point', {}).get('longitude', 0),
            'price': flat.get('price', {}).get('value', 0),
            'area': flat.get('area', {}).get('value', 0),
            'fee': flat.get('agentFee', None),
            'prepayment': flat.get('prepayment', None),
            'images': flat.get('fullImages', []),
            'planImages': flat.get('extImages', {}).get('IMAGE_PLAN', {}).get('fullImages', []),
        }

        current_data.update({
            'fee': current_data['fee'] / 100 * flat.get('price', {}).get('value', 0) if current_data['fee'] else None,
            'prepayment': current_data['prepayment'] / 100 * flat.get('price', {}).get('value', 0) if current_data['prepayment'] else None,
        })

        for image in current_data['planImages']:
            current_data['images'].remove(image)

        data.append(current_data)

    return data


def get_info(soup):
    return parse_yandex(soup)


state = requests.get(YANDEX_URL, headers=HEADERS)
soup = BeautifulSoup(state.text.encode(state.encoding).decode('utf-8'), 'html.parser')

print(json.dumps(get_info(soup), indent=2, ensure_ascii=False))