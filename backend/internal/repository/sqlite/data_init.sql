INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Foot Massage', 'I will give you a 20 minute foot massage at a time of your choosing', 1, 0, '2026-04-19 09:15:05', '2026-04-19 09:15:05');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Work Pickup', 'I will pick you up after work without a single complaint', 1, 0, '2026-04-19 09:15:30', '2026-04-19 09:15:30');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Dessert Run', 'I will go and buy whatever you are currently craving. Mango sticky rice?', 1, 1, '2026-04-19 09:16:02', '2026-04-19 09:18:45');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Dry Shower', 'I will dry the shower even if I had the first shower. Fair warning - I will be trying to shower first more so you use this one.', 1, 0, '2026-04-19 09:16:51', '2026-04-19 09:16:51');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Reveal', 'Reveals the name of an upcoming day. You may choose which day you want revealed.', 2, 1, '2026-04-19 09:17:27', '2026-04-19 09:17:27');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Jinx', 'Pick a single word I can''t say. Every time you hear me say it for the rest of the day: you gain a favour point.', 3, 1, '2026-04-19 09:17:54', '2026-04-19 09:17:54');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Baking', 'I will (attempt) to make a baked item of your choice.', 2, 1, '2026-04-19 09:18:34', '2026-04-19 09:18:34');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Cocktail Night', 'I will (attempt) to make you whatever cocktails you request', 3, 1, '2026-04-19 09:19:16', '2026-04-19 09:19:16');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Overruled', 'You can overrule something I have done today', 4, 1, '2026-04-19 09:19:45', '2026-04-19 09:19:45');
INSERT INTO favour_choice (label, description, cost, can_message, created_at, updated_at) VALUES ('Mystery', null, 7, 0, '2026-04-19 09:20:16', '2026-04-19 09:20:16');

-- Sunset Stories
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Sunset Stories', '', '2026-03-31 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:43:27', '2026-04-19 09:43:44');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Sunset Stories', '', '2026-04-31 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:43:27', '2026-04-19 09:43:44');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('dd6fb428-8f9d-4dd6-bbb8-84da91ac4320', 1, 'start', 'Start', null, null, -119.5, -25, '2026-04-19 09:43:27', '2026-04-19 21:46:58');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('23f80e93-7294-49d5-a092-c11d854752b5', 1, 'code', 'Tidal Trivia', '', null, 176.5, -95.54856965596926, '2026-04-19 21:43:30', '2026-04-19 21:46:58');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('f45133eb-6fac-4cc8-9031-ab5d67523038', 1, 'manual', 'Pack your Bags', '', null, 483.14132314716187, -79.52775084989256, '2026-04-19 21:45:18', '2026-04-19 21:46:58');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('c8d15775-ee53-4e99-9f0d-7c572c43300e', 1, 'time', 'Patience', '', null, 788.3448436456779, -94.31734875951227, '2026-04-19 21:45:18', '2026-04-19 21:46:58');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('c8d15775-ee53-4e99-9f0d-7c572c43300e', 1, 'time', 'Patience', '', null, 788.3448436456779, -94.31734875951227, '2026-04-19 21:45:18', '2026-04-19 21:46:58');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('b2715926-2910-4bd2-87bc-2729bb1b6b8f', 1, 'location', 'Follow the Sunset', '', null, 1104.3044353511898, -67.42717074202187, '2026-04-19 21:46:19', '2026-04-19 21:46:58');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('52da7429-1222-489d-bebf-f6cf05388384', 1, 'reward', 'All Aboard!', '', null, 1470.7381264090977, -170.5799841107471, '2026-04-19 21:46:19', '2026-04-19 21:46:58');
INSERT INTO edge (id, graph_id, source_node_id, destination_node_id, choice_label, created_at, updated_at) VALUES ('7a79bb84-11c5-4df6-91d9-cefd3e8cc22d', 1, 'dd6fb428-8f9d-4dd6-bbb8-84da91ac4320', '23f80e93-7294-49d5-a092-c11d854752b5', '', '2026-04-19 21:43:34', '2026-04-19 21:46:58');
INSERT INTO edge (id, graph_id, source_node_id, destination_node_id, choice_label, created_at, updated_at) VALUES ('b16f4baf-5e76-4a03-9be6-2bad6802c0cf', 1, '23f80e93-7294-49d5-a092-c11d854752b5', 'f45133eb-6fac-4cc8-9031-ab5d67523038', '', '2026-04-19 21:45:18', '2026-04-19 21:46:58');
INSERT INTO edge (id, graph_id, source_node_id, destination_node_id, choice_label, created_at, updated_at) VALUES ('71994c67-0591-48e6-833f-a8bcecfd2ed8', 1, 'f45133eb-6fac-4cc8-9031-ab5d67523038', 'c8d15775-ee53-4e99-9f0d-7c572c43300e', '', '2026-04-19 21:45:18', '2026-04-19 21:46:58');
INSERT INTO edge (id, graph_id, source_node_id, destination_node_id, choice_label, created_at, updated_at) VALUES ('8cea21da-3e3e-4a7f-8a23-fec006c7ca52', 1, 'c8d15775-ee53-4e99-9f0d-7c572c43300e', 'b2715926-2910-4bd2-87bc-2729bb1b6b8f', '', '2026-04-19 21:46:19', '2026-04-19 21:46:58');
INSERT INTO edge (id, graph_id, source_node_id, destination_node_id, choice_label, created_at, updated_at) VALUES ('e976697d-3cd9-4761-9b83-ed28f053e043', 1, 'b2715926-2910-4bd2-87bc-2729bb1b6b8f', '52da7429-1222-489d-bebf-f6cf05388384', '', '2026-04-19 21:46:19', '2026-04-19 21:46:58');

-- Waddle you do without me
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Waddle you do without me', '', '2026-04-02 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:43:59', '2026-04-19 09:44:09');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Waddle you do without me', '', '2026-05-02 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:43:59', '2026-04-19 09:44:09');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('97f399d0-0ef5-4184-8f46-92a631ecc1fb', 2, 'start', 'Start', null, null, null, null, '2026-04-19 09:43:59', '2026-04-19 09:44:09');

-- A Little Treat
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('A Little Treat', '', '2026-04-04 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:44:19', '2026-04-19 09:45:15');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('A Little Treat', '', '2026-05-04 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:44:19', '2026-04-19 09:45:15');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('f6f8b8cc-57fc-4023-af89-019fe1e0fc58', 3, 'start', 'Start', null, null, null, null, '2026-04-19 09:44:19', '2026-04-19 09:45:15');

-- Just the Beginning
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Just the Beginning', '', '2026-04-07 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:44:37', '2026-04-19 09:44:46');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Just the Beginning', '', '2026-05-07 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:44:37', '2026-04-19 09:44:46');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('b1ef5e1e-239e-4096-bef9-ca1b59f9749c', 4, 'start', 'Start', null, null, null, null, '2026-04-19 09:44:37', '2026-04-19 09:44:46');

-- Between the Lines
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Between the Lines', '', '2026-04-09 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:45:25', '2026-04-19 09:45:33');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Between the Lines', '', '2026-05-09 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:45:25', '2026-04-19 09:45:33');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('37a44805-0154-4aa9-8c50-ecfa8922fb59', 5, 'start', 'Start', null, null, null, null, '2026-04-19 09:45:25', '2026-04-19 09:45:33');

-- Worth the Wait
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Worth the Wait', '', '2026-04-10 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:45:43', '2026-04-19 09:45:52');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Worth the Wait', '', '2026-05-10 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:45:43', '2026-04-19 09:45:52');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('982fecdd-6407-44d1-b9d6-be368355d724', 6, 'start', 'Start', null, null, null, null, '2026-04-19 09:45:43', '2026-04-19 09:45:52');

-- There and Back Again
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('There and Back Again', '', '2026-04-12 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:46:10', '2026-04-19 09:46:20');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('There and Back Again', '', '2026-05-12 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:46:10', '2026-04-19 09:46:20');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('1dd132a8-4fc3-42d2-bda4-035242605c68', 7, 'start', 'Start', null, null, null, null, '2026-04-19 09:46:10', '2026-04-19 09:46:20');

-- Rough Cuts
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Rough Cuts', '', '2026-04-13 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:46:26', '2026-04-19 09:46:37');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Rough Cuts', '', '2026-05-13 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:46:26', '2026-04-19 09:46:37');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('1772f436-323b-4512-936f-d55bf18a587d', 8, 'start', 'Start', null, null, null, null, '2026-04-19 09:46:26', '2026-04-19 09:46:37');

-- Fork Me
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Fork Me', '', '2026-04-15 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:46:54', '2026-04-19 09:47:09');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Fork Me', '', '2026-05-15 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:46:54', '2026-04-19 09:47:09');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('f9b91a55-d760-4011-896e-3866a8fd3a76', 9, 'start', 'Start', null, null, null, null, '2026-04-19 09:46:54', '2026-04-19 09:47:09');

-- Grand Finale
INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Grand Finale', '', '2026-04-16 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:47:17', '2026-04-19 09:47:24');
-- INSERT INTO graph (title, description, starting_at, viewport_x, viewport_y, viewport_zoom, created_at, updated_at) VALUES ('Grand Finale', '', '2026-05-16 12:00:00+00:00', 322, 238, 2, '2026-04-19 09:47:17', '2026-04-19 09:47:24');
INSERT INTO node (id, graph_id, type, title, description, unlocked_at, ui_pos_x, ui_pos_y, created_at, updated_at) VALUES ('12987e97-b6a5-4ad1-8499-b67dbed9f0e6', 10, 'start', 'Start', null, null, null, null, '2026-04-19 09:47:17', '2026-04-19 09:47:24');

