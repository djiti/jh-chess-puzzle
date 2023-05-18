import chess

""" Run with
python3 -m pytest
"""


def test_layout_print():
    # fmt: off
    my_layout = {
                     (3, 1): 'W',
                     (2, 1): ' ', (2, 2): ' ',
                     (1, 1): ' ', (1, 2): 'W', (1, 3): ' ',
        (0, 0): 'B', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    }
    # fmt: on
    printout = chess.string_layout(my_layout)
    expected_printout = """ ♘  
 □□ 
 □♘□
♞□♞□
"""
    assert printout == expected_printout


def test_legal_move_actually_moves():
    my_board = chess.Chessboard()
    res = my_board.move_piece((0, 0), (2, 1))
    assert res
    # fmt: off
    assert my_board.layout == {
                        (3, 1): 'W',
                        (2, 1): 'B', (2, 2): ' ',
                        (1, 1): ' ', (1, 2): 'W', (1, 3): ' ',
           (0, 0): ' ', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    }
    # fmt: on


def test_move_to_busy_fails():
    my_board = chess.Chessboard()
    assert not my_board.move_piece((0, 0), (1, 2))


def test_move_from_empty_fails():
    my_board = chess.Chessboard()
    assert not my_board.move_piece((0, 1), (1, 2))


def test_move_to_outside_fails():
    my_board = chess.Chessboard()
    assert not my_board.move_piece((0, 0), (1, 2))


def test_find_black():
    my_board = chess.Chessboard()
    assert my_board.find_pieces("B") == set([(0, 0), (0, 2)])


def test_find_white():
    my_board = chess.Chessboard()
    assert my_board.find_pieces("W") == set([(1, 2), (3, 1)])


def test_find_next_boards():
    my_board = chess.Chessboard()
    next_boards = my_board.next_boards()
    # fmt: off
    EXPECTED_BOARD1 = {
                         (3, 1): 'W',
                         (2, 1): 'B', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): 'W', (1, 3): ' ',
            (0, 0): ' ', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    }
    EXPECTED_BOARD2 = {
                         (3, 1): 'W',
                         (2, 1): 'B', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): 'W', (1, 3): ' ',
            (0, 0): 'B', (0, 1): ' ', (0, 2): ' ', (0, 3): ' ',
    }
    # fmt: on
    assert len(next_boards) == 2
    assert EXPECTED_BOARD1 in [b.layout for b in next_boards]
    assert EXPECTED_BOARD2 in [b.layout for b in next_boards]
    assert [len(b.history) for b in next_boards] == [2, 2]


def test_find_next_boards_harder():
    # fmt: off
    start_layout = {
                     (3, 1): 'W',
                     (2, 1): ' ', (2, 2): ' ',
                     (1, 1): ' ', (1, 2): ' ', (1, 3): 'B',
        (0, 0): 'W', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    }
    EXPECTED_BOARDS = set([chess.string_layout(layout) for layout in [{
                         (3, 1): ' ',
                         (2, 1): ' ', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): 'W', (1, 3): 'B',
            (0, 0): 'W', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    },
                          {
                         (3, 1): 'W',
                         (2, 1): 'W', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): ' ', (1, 3): 'B',
            (0, 0): ' ', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    },
                          {
                         (3, 1): 'W',
                         (2, 1): 'B', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): ' ', (1, 3): 'B',
            (0, 0): 'W', (0, 1): ' ', (0, 2): ' ', (0, 3): ' ',
    },
                          {
                         (3, 1): 'W',
                         (2, 1): 'B', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): ' ', (1, 3): ' ',
            (0, 0): 'W', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    },
                          {
                         (3, 1): 'W',
                         (2, 1): ' ', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): ' ', (1, 3): ' ',
            (0, 0): 'W', (0, 1): 'B', (0, 2): 'B', (0, 3): ' ',
    },
                          {
                         (3, 1): 'W',
                         (2, 1): ' ', (2, 2): ' ',
                         (1, 1): ' ', (1, 2): 'W', (1, 3): 'B',
            (0, 0): ' ', (0, 1): ' ', (0, 2): 'B', (0, 3): ' ',
    },
    ]])
    # fmt: on
    my_board = chess.Chessboard(start_layout=start_layout)
    next_boards = my_board.next_boards()
    assert set([chess.string_layout(b.layout) for b in next_boards]) == EXPECTED_BOARDS


def test_easy_traverse_has_two_solutions():
    # fmt: off
    easy_start = {
                     (3, 1): 'B',
                     (2, 1): 'W', (2, 2): ' ',
                     (1, 1): ' ', (1, 2): 'B', (1, 3): ' ',
        (0, 0): ' ', (0, 1): ' ', (0, 2): 'W', (0, 3): ' ',
    }
    # fmt: on
    my_board = chess.Chessboard(start_layout=easy_start)
    completed_boards, _ = chess.traverse(my_board)
    for board in completed_boards:
        print(board)
        print(board.print_history())
    assert len(completed_boards) == 2


def test_traverse_finds_the_two_solutions():
    my_board = chess.Chessboard()
    completed_boards, _ = chess.traverse(my_board)
    assert len(completed_boards) == 2
