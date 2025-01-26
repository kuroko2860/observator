import dayjs from 'dayjs'

export const ApiStatisticDefault = {
    service_name:'',
    endpoint:'',
    method:'',
    unit:'',
    from:null,
    to:null,
}
export const DefaultFromTo = {
    from: dayjs().startOf('day').valueOf(),
    to: dayjs().startOf('day').add(1, 'day').valueOf(),
}
export const DefaultUnit = {
    unit: 'hour',
}