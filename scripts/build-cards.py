#!/usr/bin/python3

"""
This script generates the cards.json file from two given text files containing
the black and white cards.

The steps are as follows:
1. Read the black and white cards from the text files
2. For each card calculate the hash of the text
3. Construct the JSON object containing the cards
"""

from hashlib import md5
from json import dumps
from os import sep
from pathlib import Path
from sys import exit as _exit
from typing import List, Tuple


PATH = str(Path(__file__).parent.absolute())


def parse_card_content(black_file: str, white_file: str) -> Tuple[List[str], List[str]]:
    """Reads the black and white cards from the text files and returns them as a list of strings."""

    blacks: List[str]
    with open(black_file, 'r', encoding='UTF-8') as file:
        blacks = file.read().strip().split('\n')

    whites: List[str]
    with open(white_file, 'r', encoding='UTF-8') as file:
        whites = file.read().strip().split('\n')
    return blacks, whites


def build_object(blacks: List[str], whites: List[str]) -> List[object]:
    """Builds the JSON object of black and white cards."""

    cards = []
    for card in blacks:
        cards.append({'id': md5(card.encode('UTF-8')).hexdigest(), 'text': card, 'type': 0})

    for card in whites:
        cards.append({'id': md5(card.encode('UTF-8')).hexdigest(), 'text': card, 'type': 1})

    return cards


def write_object(cards: List[object], filename: str):
    """Writes the given cards in JSON format into the given file."""

    with open(filename, 'w', encoding='UTF-8') as file:
        file.write(dumps(cards, indent=4))
        file.write('\n')


if __name__ == '__main__':
    b, w = parse_card_content(sep.join([PATH, 'blacks.txt']), sep.join([PATH, 'whites.txt']))
    o = build_object(b, w)
    write_object(o, sep.join([PATH, '..', 'cards.json']))
else:
    print("This script cannot be imported")
    _exit(-1)
