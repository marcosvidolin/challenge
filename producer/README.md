# Producer

## Overview

The **Producer** is a command-line application designed to read a CSV file containing 
user data and send the user information to a message queue. 

## Features

- Read user data from a CSV file.
- Send user records to a message queue in configurable batch sizes.

## Getting Started

To run the Producer application, use the following command:

```shell
./main -f=<file.csv> -b=<batch_size>
```

## Arguments

-f=<file.csv>: Specify the path to the CSV file containing user data
-b=<batch_size>: Set the number of user records to be sent to the queue at a time
