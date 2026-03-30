-- Update receptionist password with bcrypt hash
UPDATE users 
SET password = '$2a$10$aFwigPtfBZhH31wc2lcv7.xdBwOSyxjeQ0RDtfcNTOegx2It9M/Ke'
WHERE email = 'receptionist@meetsushruta.com';

-- Verify the update
SELECT * FROM users WHERE email = 'receptionist@meetsushruta.com';
