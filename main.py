import operator
import json


def loadFont(path):
    f = open(path, encoding='utf-8')
    setting = json.load(f)
    return setting


class CardType:
    NoneCard = 0
    HighCard = 1  # 高牌
    Double_OneCard = 2  # 一对
    Double_TwoCard = 3  # 二对
    ThreeCard = 4  # 三条
    StraightCard = 6  # 顺子
    FlushCard = 7  # 同花
    GourdCard = 8  # 三条加对子（葫芦）
    FourCard = 9  # 四条
    StraightFlush = 10  # 同花顺
    RoyalFlush = 11  # 皇家同花顺


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
    "h": -1,
    "d": -2,
    "c": -3,
}

SPECIAL_STAIGHT = {2, 3, 4, 5, 14}
CARD_BIT = 16
CARD_A_PART = CARD_TABLE["A"] * 16 ** 4


class Cards:

    def __init__(self, cards):
        card_dict = dict()
        self.max = 0
        self.min = 100

        pre_color = CARD_TABLE[cards[1]]
        self.is_flush = True
        for indx in range(0, len(cards), 2):
            card_num = CARD_TABLE[cards[indx]]
            card_color = CARD_TABLE[cards[indx + 1]]

            self.max = max(self.max, card_num)
            self.min = min(self.min, card_num)
            if card_num in card_dict.keys():
                card_dict[card_num] += 1
            else:
                card_dict[card_num] = 1
            if card_color != pre_color:
                self.is_flush = False

        self.card_dict = card_dict
        self.card_type = None
        self.score = 0

    def ini_score_in_eq_case(self):
        score = 0
        sorted_card = sorted(self.card_dict.items(), key=operator.itemgetter(1, 0), reverse=True)
        for card, num in sorted_card:
            for i in range(num):
                score = score * CARD_BIT + card

        self.score += score


class BaseCompare:
    def compare(self, cards1: Cards, cards2: Cards):
        '''
        return compare result of two pairs of cards,
        if cards1 is bigger, return true, otherwise return false
        '''
        raise NotImplementedError


class FiveCardsComapre(BaseCompare):

    def compare(self, cards1: Cards, cards2: Cards):
        self.judge_card_type(cards1)
        self.judge_card_type(cards2)
        if cards1.card_type > cards2.card_type:
            return 1
        elif cards1.card_type < cards2.card_type:
            return 2
        else:
            cards1.ini_score_in_eq_case()
            cards2.ini_score_in_eq_case()
            if cards1.score > cards2.score:
                return 1
            elif cards1.score < cards2.score:
                return 2
            else:
                return 0

    def judge_card_type(self, cards: Cards):
        if len(cards.card_dict) == 5:
            self.judge_straight_type(cards)
        elif len(cards.card_dict) == 4:
            cards.card_type = CardType.Double_OneCard
        elif len(cards.card_dict) == 3:
            if max(cards.card_dict.values()) == 3:
                cards.card_type = CardType.ThreeCard
            else:
                cards.card_type = CardType.Double_TwoCard
        elif len(cards.card_dict) == 2:
            if max(cards.card_dict.values()) == 4:
                cards.card_type = CardType.FourCard
            else:
                cards.card_type = CardType.GourdCard
        if cards.is_flush:
            cards.card_type = max(cards.card_type, CardType.FlushCard)

    def judge_straight_type(self, cards: Cards):
        if not self.base_judge_staight(cards):
            cards.card_type = CardType.HighCard
        elif not cards.is_flush:
            cards.card_type = CardType.StraightCard
        elif cards.max == CARD_TABLE['A']:
            cards.card_type = CardType.RoyalFlush
        else:
            cards.card_type = CardType.StraightFlush

    def base_judge_staight(self, cards: Cards):
        if cards.max - cards.min == 4:
            return True
        if SPECIAL_STAIGHT.issubset(cards.card_dict.keys()):
            cards.score = (cards.score - CARD_A_PART) * CARD_BIT + 1
            cards.max = 0
            return True
        return False


if __name__ == "__main__":

    t = loadFont("test_file/result.json")["matches"]
    compartor: BaseCompare
    compartor = FiveCardsComapre()

    for i, game in enumerate(t):
        alice_card = Cards(game["alice"])
        bob_card = Cards(game["bob"])
        res = compartor.compare(alice_card, bob_card)
        assert(game['result'] == res)
        # print(i, end=" ")

    print("Are you happy?")
