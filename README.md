# go-life

A go implementation of Conway's game of life. The program takes input from stdin. It's recommended to use it as `cat input.txt | go-life` with input.txt being:

```
(period in ms)
(width of the window) (height of the window)
(number of columns) (number of rows)
(coordinates of the upper-left coordinate of the patch, column and row)
(the patch)

```

Example of input.txt:

```
0.5
500 500
500 500
250 250
.........
..+......
....+....
.++..+++.
.........

```
