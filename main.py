from enum import Enum, unique
import inp


@unique
class CardType(Enum):
    NoneCard = 0,
    HighCard = 1,  # 高牌
    Double_OneCard = 2,  # 一对
    Double_TwoCard = 3,  # 二对
    ThreeCard = 4,  # 三条
    StraightCard = 6,  # 顺子
    FlushCard = 7,  # 同花
    GourdCard = 8,  # 三条加对子（葫芦）
    FourCard = 9,  # 四条
    StraightFlush = 10,  # 同花顺
    RoyalFlush = 11  # 皇家同花顺


scores = [
    1000,
    500,
    100,
]

CARD_TABLE = {
    "2": 2,
    "3": 3,
    "4": 4,
    "5": 5,
    "6": 6,
    "7": 7,
    "8": 8,
    "9": 9,
    "T": 10,
    "J": 11,
    "Q": 12,
    "K": 13,
    "A": 14,
    "s": 0,
    "h": 1,
    "d": 2,
    "c": 3,
}


class Cards:

    def __init__(self):
        pass



class BaseCompare:
    def compare(self, cards1, cards2):
        '''
        return compare result of two pairs of cards, 
        if cards1 is bigger, return true, otherwise return false
        '''
        raise NotImplementedError


class FiveCardsComapre(BaseCompare):
    def compare(self, cards1, cards2):
        pass


if __name__ == "__main__":

    t = inp.loadFont("poke/five_cards_with_ghost.json")["matches"]
    compartor: BaseCompare
    compartor = FiveCardsComapre()

    for game in t[0:10]:
        print(compartor.compare(game["alice"], game["bob"]))

    print("Are you happy")
