ALTER TABLE Alert 
    ADD fingerprint TEXT NOT NULL default ''
;

UPDATE Model  SET version='0.1.0';