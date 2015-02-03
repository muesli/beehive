Hives.HivesController = Ember.ArrayController.extend({
  actions: {
    createHive: function() {
      var name = this.get('newName');
      if (!name.trim()) { return; }

      var hive = this.store.createRecord('hive', {
        name: name,
        isCompleted: false
      });

      this.set('newName', '');

      // Save the new model
      hive.save();
    }
  },
  remaining: function() {
    return this.filterBy('isCompleted', false).get('length');
  }.property('@each.isCompleted'),

  inflection: function() {
    var remaining = this.get('remaining');
    return remaining === 1 ? 'item' : 'items';
  }.property('remaining')
});
