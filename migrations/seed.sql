-- Create the admin user.
-- If used on production, this user's password must be modified as soon as
-- possible.
INSERT INTO users(
	username,
	password,
	email,
	first_name,
	last_name,
	status,
	created_at,
	updated_at
) VALUES (
	'admin',
	'$2a$10$vWQ6Tsu8odmvJ74.dlSZT.XonFRKuqxd/bZzxHP041FOWTzhKj552',
	'admin@abemar.com',
	'Administrador',
	'Administrador',
	0,
	timezone('UTC', now()),
	timezone('UTC', now())
)
