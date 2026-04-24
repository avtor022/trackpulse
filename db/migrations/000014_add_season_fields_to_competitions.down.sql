-- Remove season_name and competition_year columns from competitions table
ALTER TABLE competitions DROP COLUMN season_name;
ALTER TABLE competitions DROP COLUMN competition_year;
