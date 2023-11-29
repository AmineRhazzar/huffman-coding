package main

const __DEBUG__ = false

/*
				20
	    -----------------
	    12              8
	-----------     ---------
	C,6     A,6     E,5     3
						-----------
						D,2     B,1

encoded as 00101000 01110100 00010101 00010101 01000100 10100001 00000000 00000001
*/
var MOCK_TREE *Node = &Node{
	weight: 20,
	Left: &Node{
		weight: 12,
		Left:   &Node{ch: 'C', weight: 6},
		Right:  &Node{ch: 'A', weight: 6},
	},
	Right: &Node{
		weight: 8,
		Left:   &Node{ch: 'E', weight: 5},
		Right: &Node{weight: 3,
			Left:  &Node{ch: 'D', weight: 2},
			Right: &Node{ch: 'B', weight: 1},
		},
	},
}

const MOCK_TREE_ENCODED_SIZE int = 49

/*
                        23
		----------------------------------
		12                              11
	------------                -------------------
    A,6      C,6                6               E,5
						-----------------
						F,3             3
									-----------
									D,2     B,1

encoded as 00101000 01110100 00010010 10001100 10100010 01010000 10101000 10100000 00000011
*/
var MOCK_TREE_2 *Node = &Node{
	weight: 23,
	Left: &Node{
		weight: 12,
		Left:   &Node{ch: 'C', weight: 6},
		Right:  &Node{ch: 'A', weight: 6},
	},
	Right: &Node{
		weight: 11,
		Left:   &Node{ch: 'E', weight: 5},
		Right: &Node{weight: 6,
			Left: &Node{ch: 'F', weight: 3},
			Right: &Node{
				weight: 3,
				Left:   &Node{ch: 'D', weight: 2},
				Right:  &Node{ch: 'B', weight: 1},
			},
		},
	},
}

const MOCK_TREE_2_ENCODED_SIZE int = 59
