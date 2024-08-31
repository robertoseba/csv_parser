# CSV Parser Tool

![csv parser interface](docs/screen.gif)

A simple tool for parsing CSV files in the CLI with options for column and row filtering.

Here's an example of how you can use this tool:

Sometimes you got a big CSV file and want to only get rows for a single user.

```bash
csv_parser user-list.csv --filter "username, email" --rules "email:eq(user1@example.com)"
```

The command above will print only the columns username and email and only the rows that email match "user1@example.com"

## Flags:

- <b>filter</b> (optional) - Applies filtering to column names.
- <b>rules</b> (optional) - Applies filtering to rows based on rules defined for each column.<u> It works independently from the filter flag.</u> For instance, you can filter to show only a user's name and apply a rule to only get rows by the user's email.

## Filter

It's based on the header values and takes the orders of the columns as inputed in the cli command.
It can also be useful for when you want to create a new CSV file based on the original one but only want certain columns.

Let's say I have a csv file (users.csv) with the following headers: `name,email,phone,address`.

If I want to generate a new file with `email,name` in that order, here's the cli command for it:

```bash
csv_parser users.csv --filter "email" > newfile.csv
```

## Rules

You can apply multiple rules to a column. Rules for each column must be separated by `;`

Syntax: `<column-name>:<rule-type>(<value>)<optional-logical-operator><rule-type>(<value>);`

### Example:

```bash

$ csv_parser --rules "col1:eq(bob)||eq(junior);col2:neq(10)&&lte(20)"

```

> If the rule value is a number the parser will try to compare values as being numerical. Otherwise values will be compared lexicographically.

### Rule types:

- eq - Equal
- neq - Not equals
- lt - Lower than
- lte - Lower than or equal to
- gt - Greater than
- gte - Greater than or equal to

### Logical Operators:

- && - AND
- || - OR

### Use case example:

Given the following csv (users.csv):

```csv
name,score,test
bob,100,test_1
bob,200,test_2
bob,10,test_3
junior,20,test_2
junior,25,test_1
junior,100,test_3
mike,50,test_1
mike,5,test_2
mike,100,test_3
```

I want to retrieve only names for rows where test_2 or test_3 score is 100 or more.

```bash
csv_parser users.csv --filter "name" --rules "test:eq(test_2)||eq(test_3);score:gte(100)"
```
