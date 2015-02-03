Hives.Router.map(function() {
  this.resource('hives', { path: '/' });
});

Hives.HivesRoute = Ember.Route.extend({
  model: function(params) {
//    return this.store.find('hive');
    return this.store.find('hive');//, params.hive_id);
  }
});