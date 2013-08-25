`sptt`
======

Split data into a Percentage of Test and Train sets.

To split the file named 'data' into 'data.train', containing ~90% training data and another file named 'data.test', containing ~10% testing data:

    $ sptt -train 90 data

Usage for files, where `PERCENT` is an integer and `FILE` is a relative path:

    sptt -train PERCENT FILE

Usage for STDIN (this writes 2 files; `STDIN.test` and `STDIN.train`):

    sptt -train PERCENT -
