import { configureStore } from '@reduxjs/toolkit';
import zonesReducer from '../features/zones/zonesSlices';

export default configureStore({
  reducer: {
    zones: zonesReducer,
  },
});
