--@ddl:start
CREATE OR REPLACE FUNCTION a_function()
RETURNS BOOLEAN AS $$
    SELECT true;
$$ LANGUAGE sql STABLE; 
--@ddl:end

--@ddl:start
DROP POLICY IF EXISTS a_policy
ON employees;
CREATE POLICY a_policy
ON employees
FOR SELECT
TO postgres;
--@ddl:end


DO
$do$
BEGIN
  --@ddl:start
  IF EXISTS (
    SELECT FROM pg_catalog.pg_roles
    WHERE  rolname = 'postgres') THEN

    RAISE NOTICE 'Role "postgres" already exists. Skipping.'
  ELSE
    EXECUTE format('CREATE USER postgres WITH PASSWORD %L', 'postgres')
  END IF;
  --@ddl:end
END
$do$;