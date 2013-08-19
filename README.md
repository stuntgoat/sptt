`asprl`
======

Approximately Split a file into percentages of Randomly chosen Lines.

This is used for dividing a dataset into testing and training and validation sets.

Split a file into a random Percentage of Testing and Training data files.

Usage:

    # split the file named 'data' into 'data.train', containing 10% training data
    # and another file named 'data.test'
    $ sptt -train 90 data

    # another way to accomplish the same task
    $ sptt -test 10 data

    # create 4 files that divides a file between 4 different files of random lines
    $ sptt -validation 4 data
