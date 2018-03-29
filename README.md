Just Expense It
===============

This utility is a simple Expensify API client which allows you to see how much you've expensed in a certain time period.  
Currently it can only do this but it is sure expandable to more cases in the future. 

## API Client
This repository also includes an Expensify API client in Go which is heavily abstracted to provide a modern way of accessing the information in Expensify. 

## Usage
The main program currently accepts 2 environment variables `JUSTEXPENSEIT_EXPPENSIFY_USER_ID` and `JUSTEXPENSEIT_EXPPENSIFY_USER_SECRET`, these are used to authenticate against Expensify (you can get them from https://www.expensify.com/tools/integrations/). They can also be replaced by a config file (we use viper) or CLI flags.
