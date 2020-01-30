from typing import List


class Solution:
    def gameOfLife(self, board: List[List[int]]) -> None:
        """
        Do not return anything, modify board in-place instead.
        """
        yl = len(board)
        if yl == 0:
            return
        xl = len(board[0])

        for y in range(yl):
            for x in range(xl):

                # Count
                c = 0
                for i in range(-1, 2):
                    for j in range(-1, 2):
                        if not (j == 0 and i == 0):
                            if 0 <= y + i < yl and 0 <= x + j < xl:
                                c += board[y + i][x + j] % 2

                # Transition
                if board[y][x] in [1, 2]:
                    board[y][x] = 3 if c < 2 or c > 3 else 1
                else:
                    board[y][x] = 2 if c == 3 else 0

        for y in range(yl):
            for x in range(xl):
                board[y][x] = 1 if board[y][x] in [1, 2] else 0
