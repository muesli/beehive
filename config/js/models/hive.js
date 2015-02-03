Hives.Hive = DS.Model.extend({
  name: DS.attr('string'),
  image: DS.attr('string'),
  description: DS.attr('string'),
  isCompleted: DS.attr('boolean')
});

Hives.HiveAdapter = DS.RESTAdapter.extend({
  namespace: 'v1',
  host: 'http://localhost:8181'
});
