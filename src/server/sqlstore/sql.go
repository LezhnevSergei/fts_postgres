package sqlstore

const SaveTSVector = `
insert into incident_tvectors (incident_id, tvs)
select incident_id, make_tsvector(r.display_name, r.description)
from incidents
         left join rules r on incidents.rule_id = r.rule_id;
`

const ListRules = `
select r.display_name,
       r.description
from rules as r;
`

const SearchIncidents = `
select inc.incident_id,
       r.display_name,
       r.description
from plainto_tsquery($1) as q, incidents as inc
         left join rules as r on inc.rule_id = r.rule_id
         left join rules_tvectors as r_ts on r.rule_id = r_ts.rule_id
where r_ts.tvs @@ q
order by ts_rank(r_ts.tvs, q) desc
limit 50;
`
