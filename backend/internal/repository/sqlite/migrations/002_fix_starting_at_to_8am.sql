-- Fix graphs that were created with midnight UTC starting_at due to a bug
-- where date.toISOString() was used instead of the 8am-adjusted date.
-- This shifts any starting_at stored at exactly midnight UTC to 08:00 UTC.
UPDATE graph
SET starting_at = datetime(starting_at, '+8 hours')
WHERE time(starting_at) = '00:00:00';
