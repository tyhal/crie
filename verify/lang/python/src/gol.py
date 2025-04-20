"""A Simple Inplace Game of Life"""


class Solution:
    """Store board state and simulate game of life"""

    def __init__(self, board_):
        self.board = board_
        self.ylen = len(board_)
        if self.ylen == 0:
            return
        self.xlen = len(board_[0])

    def neigh(self, ypos, xpos):
        """
        :param ypos: the current y position
        :param xpos: the current x position
        :return: the number of neighbors
        """
        count = 0
        for i in range(-1, 2):
            for j in range(-1, 2):
                if (
                    not (j == 0 and i == 0)
                    and 0 <= ypos + i < self.ylen
                    and 0 <= xpos + j < self.xlen
                ):
                    count += self.board[ypos + i][xpos + j] % 2
        return count

    def iter(self) -> None:
        """
        Do not return anything, modify self.board in-place instead.
        """
        for ypos in range(self.ylen):
            for xpos in range(self.xlen):

                # Count
                count = self.neigh(ypos, xpos)

                # Transition
                if self.board[ypos][xpos] in [1, 2]:
                    self.board[ypos][xpos] = 3 if count < 2 or count > 3 else 1
                else:
                    self.board[ypos][xpos] = 2 if count == 3 else 0

        for ypos in range(self.ylen):
            for xpos in range(self.xlen):
                self.board[ypos][xpos] = 1 if self.board[ypos][xpos] in [1, 2] else 0


if __name__ == "__main__":
    INP = [[0, 0, 0], [0, 0, 1], [1, 1, 1], [0, 0, 0]]
    GOL = Solution(INP)
    for _ in range(5):
        GOL.iter()
    for row in INP:
        print(row)
