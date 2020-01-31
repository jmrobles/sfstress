#!/bin/bash
. creds.sh
# echo $SF_DB
echo "# Query 1: SELECT * con IDs menores 5.000.000"
echo "# -- ODBC Snowflake"
./sfstress.exe -backend odbc -bulkFile sql/001_query.sql -mode bulk
./sfstress.exe -backend odbc -bulkFile sql/001_query.sql -mode bulk
./sfstress.exe -backend odbc -bulkFile sql/001_query.sql -mode bulk

echo "# -- JDBC Snowflake"
./sfstress.exe -backend snowflake -bulkFile sql/001_query.sql -mode bulk
./sfstress.exe -backend snowflake -bulkFile sql/001_query.sql -mode bulk
./sfstress.exe -backend snowflake -bulkFile sql/001_query.sql -mode bulk

echo "# -- SQL Server"
./sfstress.exe -backend sqlserver -bulkFile sql/001_query_mssql.sql -mode bulk
./sfstress.exe -backend sqlserver -bulkFile sql/001_query_mssql.sql -mode bulk
./sfstress.exe -backend sqlserver -bulkFile sql/001_query_mssql.sql -mode bulk

echo "# Query 2: SUMA DEL DÉBITO"
echo "# -- ODBC Snowflake"
./sfstress.exe -backend odbc -bulkFile sql/002_sumcnt.sql -mode bulk
./sfstress.exe -backend odbc -bulkFile sql/002_sumcnt.sql -mode bulk
./sfstress.exe -backend odbc -bulkFile sql/002_sumcnt.sql -mode bulk
echo "# -- JDBC Snowflake"
./sfstress.exe -backend snowflake -bulkFile sql/002_sumcnt.sql -mode bulk
./sfstress.exe -backend snowflake -bulkFile sql/002_sumcnt.sql -mode bulk
./sfstress.exe -backend snowflake -bulkFile sql/002_sumcnt.sql -mode bulk
echo "# -- SQL Server"
./sfstress.exe -backend sqlserver -bulkFile sql/002_sumcnt.sql -mode bulk
./sfstress.exe -backend sqlserver -bulkFile sql/002_sumcnt.sql -mode bulk
./sfstress.exe -backend sqlserver -bulkFile sql/002_sumcnt.sql -mode bulk
