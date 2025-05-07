-- Drop product-specific schemas

-- Drop Invantray schema
DROP TABLE IF EXISTS invantray.inventory_transactions;
DROP TABLE IF EXISTS invantray.warehouses;
DROP TABLE IF EXISTS invantray.inventory_items;
DROP SCHEMA IF EXISTS invantray;

-- Drop CustomerConnect schema
DROP TABLE IF EXISTS customerconnect.opportunities;
DROP TABLE IF EXISTS customerconnect.leads;
DROP TABLE IF EXISTS customerconnect.contacts;
DROP SCHEMA IF EXISTS customerconnect;

-- Drop AdminHub schema
DROP TABLE IF EXISTS adminhub.system_settings;
DROP TABLE IF EXISTS adminhub.system_logs;
DROP SCHEMA IF EXISTS adminhub;

-- Drop Qultrix schema
DROP TABLE IF EXISTS qultrix.performance_reviews;
DROP TABLE IF EXISTS qultrix.time_off_requests;
DROP TABLE IF EXISTS qultrix.employees;
DROP SCHEMA IF EXISTS qultrix;

-- Drop Taskodex schema
DROP TABLE IF EXISTS taskodex.task_comments;
DROP TABLE IF EXISTS taskodex.tasks;
DROP TABLE IF EXISTS taskodex.projects;
DROP SCHEMA IF EXISTS taskodex;
