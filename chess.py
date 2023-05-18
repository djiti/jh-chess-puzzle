#!/usr/bin/env python3

# Knights will be modeled as 'W' and 'B' for their color
# Using enums would be cleaner, but hey

from copy import deepcopy, copy


PIECE_MAPPING = {"W": "♘", "B": "♞", " ": "□"}


def string_layout(layout):
    """
    Returns  an ASCII representation of a layout, handy for "hashing"
    (as our layouts are dictionaries and hence not hashable)
    and for printing.
    Layouts are indexed by their column and row number,
    starting on the bottom left corner
    ' ' if empty, 'W' if White knight and 'B' if Black knight"""
    printout = ""
    for row in range(3, -1, -1):
        for col in range(4):
            position = (row, col)
            if position in layout:
                slot = PIECE_MAPPING[layout[position]]
            else:
                slot = " "
            printout = printout + slot
        printout = printout + "\n"
    return printout


class Chessboard:
    # fmt: off
    START_LAYOUT = {
        (0, 0): 'B', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
                     (1, 1): ' ', (1, 2): 'W', (1, 3): ' ',
                     (2, 1): ' ', (2, 2): ' ',
                     (3, 1): 'W'
    }
    GOAL_LAYOUT = {
        (0, 0): 'W', (0, 1): ' ', (0, 2): 'W', (0, 3): ' ',
                     (1, 1): ' ', (1, 2): 'B', (1, 3): ' ',
                     (2, 1): ' ', (2, 2): ' ',
                     (3, 1): 'B'
    }
    # fmt: on
    def __init__(self, start_layout=START_LAYOUT, goal=GOAL_LAYOUT):
        """A Chessboard is a layout, mechanisms to mutate its layout and
        a place to keep its history."""
        self.layout = copy(start_layout)
        self.goal_layout = copy(goal)
        self.history = [copy(self.layout)]

    @property
    def finished(self):
        """Property when the board has reached the desired layout"""
        return self.layout == self.goal_layout

    def find_pieces(self, piece):
        """Gets the list of pieces"""
        return set(
            [value for value in self.layout.keys() if self.layout[value] == piece]
        )

    def print_history(self):
        """Returns the board's full history"""
        a = ""
        for the_layout in self.history:
            a = f"{a}\n{string_layout(the_layout)}"
        return a

    def move_piece(self, start_pos, end_pos):
        """Moves a piece while keeping a history of all moves.
        Returns True in case of success, False in case of failure.
        """
        if start_pos not in self.layout or end_pos not in self.layout:
            return False
        if self.layout[end_pos] != " ":
            return False
        if self.layout[start_pos] == " ":
            return False
        piece = self.layout[start_pos]
        self.layout[start_pos] = " "
        self.layout[end_pos] = piece
        self.history.append(deepcopy(self.layout))
        return True

    def next_boards(self):
        """Constructor for the possible next boards"""
        boards = []
        for start_pos in self.find_pieces("B").union(self.find_pieces("W")):
            # fmt: off
            for move in [(2, -1), (2, 1), (1, -2), (1, 2), (-2, -1), (-2, 1), (-1, 2), (-1, -2)]:
                # fmt: on
                new_board = deepcopy(self)
                end_pos = (start_pos[0]+move[0], start_pos[1]+move[1])
                if new_board.move_piece(start_pos, end_pos):
                    boards.append(new_board)
        return boards


def traverse(board):
    """Traverses all possible boards and returns the boards that reach the goal layout
    and the layouts it visited along the way (excluding the final layout...)"""
    finished_boards = []
    boards = [board]
    visited_layouts = {}
    while boards:
        visited_layouts.update({string_layout(b.layout): b.layout for b in boards})
        next_boards = [
            inner_board
            for outer_board in boards
            for inner_board in outer_board.next_boards()
            # This takes care of removing both
            # - looping boards (they are not interesting if they are not
            # finished by now)
            # - layouts that have already been visited
            if not inner_board.layout in visited_layouts.values()
        ]
        # We could decide to stop at the first finished board, but let's find more
        # (all?) of them instead
        finished_boards.extend(
            [the_board for the_board in next_boards if the_board.finished]
        )
        # We keep only one board per current layout
        # (so technically we throw the others' history away,
        # which may ignore finished boards - to be pondered)
        # and exclude all those that have reached the goal
        boards = {
            string_layout(the_board.layout): the_board
            for the_board in next_boards
            if not the_board.finished
        }.values()

    return finished_boards, visited_layouts.values()


def main():
    board = Chessboard()
    completed_boards, visited_layouts = traverse(board)
    for board in completed_boards:
        print("*********************")
        print(
            f"In {len(board.history)} moves, "
            f"visiting {len(visited_layouts)} layouts (excluding the final one):"
        )
        print(board.print_history())
        print()


if __name__ == "__main__":
    main()
