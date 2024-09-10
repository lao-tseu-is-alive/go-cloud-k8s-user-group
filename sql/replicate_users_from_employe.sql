with old as (select upper(nom) ||  ', ' ||  prenom                                 as name,
                    email,
                    mainntlogin                                                    as username,
                    '$2a$04$Aj1Frkj5lvlf3WWTQIUqB.WF7SqYpFC7ZnoI3/OA5xpcHyoDL6jv6' as password_hash,
                    idemploye                                                      as external_id,
                    telprof                                                        as phone,
                    datecreated                                                    as create_time,
                    idcreator                                                      as creator,
                    datelastmodif                                                  as last_modification_time,
                    idlastmodifuser                                                as last_modification_user,
                    isactive                                                       as is_active,
                    comment                                                        as comment,
                    0                                                              as bad_password_count
             from employe
             where isactive = true
               AND employe.mainntlogin NOT like '%BIS%')
INSERT
INTO go_user (name, email, username, password_hash, external_id, phone, create_time, creator, last_modification_time,
              last_modification_user, is_active, comment, bad_password_count)
SELECT *
FROM old;
