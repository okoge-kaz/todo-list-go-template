INSERT INTO `tasks` (`title`) VALUES ("sample-task-01");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-02");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-03");
INSERT INTO `tasks` (`title`, `description`) VALUES ("sample-task-04", "sample-task-04 description");
INSERT INTO `tasks` (`title`, `is_done`) VALUES ("sample-task-05", true);

INSERT INTO `users` (`name`, `password`) VALUES ("kazuki", "password");
INSERT INTO `users` (`name`, `password`) VALUES ("okoge", "pass");

INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 1);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 2);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 3);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (2, 4);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (2, 5);
