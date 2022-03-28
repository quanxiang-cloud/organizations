alter table org_user change user_name name varchar (64) null;
alter table org_user change create_time created_at bigint null;
alter table org_user change update_time updated_at bigint null;
alter table org_user change create_by created_by varchar (64) null;
alter table org_user change update_by updated_by varchar (64) null;
alter table org_user
    add deleted_at bigint null;
alter table org_user
    add deleted_by bigint null;
alter table org_user drop column company_id;
alter table org_user drop column leader_id;


alter table org_user_account
    add column account varchar(100) null;
alter table org_user_account
    add column user_id varchar(100) null;
alter table org_user_account
    add column tenant_id varchar(100) null;
alter table org_user_account change create_time created_at bigint null;
alter table org_user_account change update_time updated_at bigint null;
alter table org_user_account change create_by created_by varchar (64) null;
alter table org_user_account change update_by updated_by varchar (64) null;
alter table org_user_account
    add deleted_at bigint null;
alter table org_user_account
    add deleted_by bigint null;



alter table org_department change creat_by create_by varchar (64) null;
alter table org_department change create_time create_at bigint null;
alter table org_department change update_time update_at bigint null;
alter table org_department
    add delete_at bigint null;
alter table org_department
    add delete_by bigint null;
alter table org_department drop column leader_id;
alter table org_department drop column company_id;
alter table org_department drop column third_id;
