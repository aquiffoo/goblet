create table users (
    id integer primary key autoincrement,
    username text not null,
    password text not null
);

create table todos (
    id integer primary key autoincrement,
    title text not null,
    description text not null,
    done boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    user_id integer not null,
    foreign key (user_id) references users(id)
);
