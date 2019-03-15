import mongoose from 'mongoose';

const DataSchema = mongoose.Schema({
  from: {
    type: String,
    required: true,
  },
  timestamp: {
    type: Date,
    default: Date.now,
    required: true,
  },
  payload: {
    type: Object,
    required: true,
  },
});

class DataStore {
  constructor(database) {
    this.database = database;
  }

  async save(data) {
    await this.database.save('Data', DataSchema, data);
  }
}

export default DataStore;
