#!/bin/bash

rm -rf github

buf -v generate proto
