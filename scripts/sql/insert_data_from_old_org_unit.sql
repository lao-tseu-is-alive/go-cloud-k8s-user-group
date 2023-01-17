INSERT INTO public.go_orgunit
(id, type, name, parent_id, abbreviation, description, order_list, phone, email, create_time, creator, comment, guid,
 full_name_de_norm)
SELECT id,
       type,
       name,
       parent_id,
       abbreviation,
       description,
       orderlist,
       phone,
       email,
       create_time,
       creator,
       comment,
       guid,
       desctreedenorm
FROM org_unit
WHERE org_unit.type IN ('Entreprise', 'Direction', 'Service')
ORDER BY orderlist;

-- insert one level of ou above service
INSERT INTO public.go_orgunit
(id, type, name, parent_id, abbreviation, description, order_list, phone, email, create_time, creator, comment, guid,
 full_name_de_norm)
SELECT id,
       type,
       name,--length(name),
       parent_id,
       abbreviation,
       description,
       orderlist,
       phone,
       email,
       create_time,
       creator,
       comment,
       guid,
       desctreedenorm
FROM org_unit
WHERE parent_id IN (SELECT id FROM org_unit WHERE type_orgunit_id = 3 ORDER BY 1)
ORDER BY 4;

SELECT id, type_orgunit_id, name, abbreviation, orderlist, desctreedenorm
FROM org_unit
WHERE type_orgunit_id = 3
ORDER BY orderlist

--DELETE FROM org_unit WHERE id=24

SELECT abbreviation, count(*)
FROM org_unit
WHERE parent_id IN (SELECT id FROM org_unit WHERE type_orgunit_id = 3 ORDER BY 1)
GROUP BY abbreviation
ORDER BY 2 DESC
