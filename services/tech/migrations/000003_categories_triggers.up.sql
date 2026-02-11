create or replace function null_category() returns trigger as
$func$
begin
execute format('update technics set category_id = uuid_nil() where category_id=$1') using old.id;
return old; 
end
$func$ language plpgsql;

create or replace trigger null_category_on_delete before delete on categories for each row execute procedure null_category();