# Inverse Distance Trees implemented in Go

This is a partial implementation, and will be replaced if I can find another IDT implementation out there.

Proportional and majority algorithms are both implemented.

## Tests

The input grid has a upper-left coordinate of (0, 0) and a spacing of
8 in each direction.

0 | 0 | 1 | 1
--+---+---+--
0 | 1 | 1 | 1
--+---+---+--
0 | 1 | 2 | 3
--+---+---+--
1 | 2 | 3 | 3

The output grid has an upper-left coordinate of (2, 2) and a spacing of 4 in each direction.

The four nearest neighbors to the output grid point (22, 14) are:

  Point  | Distance | Inverse |   Weight
---------+----------+---------+-----------
(24, 16) |    8     |  0.125  | 0.66176471
(16, 16) |   40     |  0.025  | 0.13235294
(24,  8) |   40     |  0.025  | 0.13235294
(16,  8) |   72     |  0.014  | 0.07352941
---------+----------+---------+-----------
  Total  |  160     |  0.189  | 1.00000000

The points are weighed by the inverse of their distance.  First, the distances are inverted and summed.  Then the inverses are divided by the sum of the inverses so they are normalized to sum to 1.0.  These weights are used for both proportional and majority algorithms.

### Proportional

Using 4 nearest neighbors, the result should be:

 0 | 0 | 0 | 1 | 1 | 1
 --+---+---+---+---+--
 0 | 1 | 1 | 1 | 1 | 1
 --+---+---+---+---+--
 0 | 1 | 1 | 1 | 1 | 1
 --+---+---+---+---+--
 0 | 1 | 1 | 2 | 2 | 2
 --+---+---+---+---+--
 0 | 1 | 1 | 2 | 2 | 3
 --+---+---+---+---+--
 1 | 2 | 2 | 3 | 3 | 3

In the proportional algorithm, the weights are multiplied by the
values and the sum of the results represents the output value.  If the
result is not an integer, standard rounding rules apply.

Using the data for the output grid point (22,14):

  Point  | Value |  Portion
---------+-------+-----------
(24, 16) |   3   | 1.98529413
(16, 16) |   2   | 0.26470588
(24,  8) |   1   | 0.13235294
(16,  8) |   1   | 0.07352941
---------+-------+-----------
  Total          | 2.45588236

In this case, the value for that point would be 2.

### Majority

Using the majority algorithm, the result should be:

 0 | 0 | 0 | 1 | 1 | 1
 --+---+---+---+---+--
 0 | 1 | 1 | 1 | 1 | 1
 --+---+---+---+---+--
 0 | 1 | 1 | 1 | 1 | 1
 --+---+---+---+---+--
 0 | 1 | 1 | 2 | 2 | 3
 --+---+---+---+---+--
 0 | 1 | 1 | 2 | 2 | 3
 --+---+---+---+---+--
 1 | 2 | 2 | 3 | 3 | 3

In the majority algorithm, the weights of each value found within the
neighbors are summed, and the value with the largest sum is returned.

Using the data for the output grid point (22,14):

 Value |   Total    | Source
-------+------------+--------
   3   | 0.66176471 | (24, 16)
   2   | 0.13235294 | (16, 16)
   1   | 0.20588235 | (24, 8) and (16, 8)
-------+------------+--------
 Total | 1.00000000 |
 
In this case, the value for that point would be 3.




