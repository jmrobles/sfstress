#!/bin/bash
. creds.sh

./sfstress.exe -backend odbc -bulkFile sql/001_query.sql -mode bulk
# ./sfstress.exe -backend odbc -bulkFile sql/001_query.sql -mode bulk
# ./sfstress.exe -backend odbc -bulkFile sql/001_query.sql -mode bulk

# ./sfstress.exe -backend snowflake -bulkFile sql/001_query.sql -mode bulk
# ./sfstress.exe -backend snowflake -bulkFile sql/001_query.sql -mode bulk
# ./sfstress.exe -backend snowflake -bulkFile sql/001_query.sql -mode bulk

# ./sfstress.exe -backend sqlserver -bulkFile sql/001_query_mssql.sql -mode bulk
# ./sfstress.exe -backend sqlserver -bulkFile sql/001_query_mssql.sql -mode bulk
# ./sfstress.exe -backend sqlserver -bulkFile sql/001_query_mssql.sql -mode bulk

