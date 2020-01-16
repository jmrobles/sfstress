# Snowflake Stress Tester

This simple tool executes concurrent queries to stress Snowflake & SQL Server databases.

## Usage

 ```bash
sfstress <sql-test-path> <backend> -duration 300 -concurrent 10
 ```
 Where:

* __sql-test-path__, path where the SQL test files to run are located
* __backend__, "snowflake" for Snowflake Database or "sqlserver" for Microsoft SQL Server
* __duration__, time in seconds for total execution
* __concurrent__, number of parallel queries threads

JM Robles @ Civica Software &copy; 2019

MIT Licence
