import {configureStore} from '@reduxjs/toolkit'
import servicesReducer from './services/slice'

export default configureStore({
    reducer:{
        services: servicesReducer
    }
})