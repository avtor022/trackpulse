-- Add season_name and competition_year columns to competitions table
ALTER TABLE competitions ADD COLUMN season_name TEXT;
ALTER TABLE competitions ADD COLUMN competition_year INTEGER;
