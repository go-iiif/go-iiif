Mean
====

A simple color quantizer.  Similar to a median cut algorithm execept it uses
the mean rather than the median.  While the median seems technically correct
the mean seems good enough and is easier to compute.  This implementation is
also simplified over traditional median cut algorithms by replacing the
priority queue with a simple linear search.  For a typical number of colors
(256) a linear search is fast enough and does not represent a significant
inefficiency.

A bit of sophistication added though is a two stage clustering process.
The first stage takes a stab at clipping tails off the distributions of colors
by pixel population, with the goal of smoother transitions in larger areas of
lower color density.  The second stage attempts to allocate remaining palette
entries more uniformly.  It prioritizes by a combination of pixel population
and color range and splits clusters at mean values.
