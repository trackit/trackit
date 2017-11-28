//import { call } from './misc.js';

const mockData = [
  {
    "_id": "bucket-0",
    "size": 71763248099,
    "storage_cost": 71.76324809900001,
    "bw_cost": 262.01,
    "total_cost": 333.773248099,
    "transfer_in": 192264344,
    "transfer_out": 134694718,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-1",
    "size": 347987615064,
    "storage_cost": 347.987615064,
    "bw_cost": 137.1,
    "total_cost": 485.08761506400003,
    "transfer_in": 270770700,
    "transfer_out": 358599956,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-2",
    "size": 139981408237,
    "storage_cost": 139.981408237,
    "bw_cost": 770.75,
    "total_cost": 910.731408237,
    "transfer_in": 78680820,
    "transfer_out": 326039295,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-3",
    "size": 329958421801,
    "storage_cost": 329.95842180100004,
    "bw_cost": 950.36,
    "total_cost": 1280.318421801,
    "transfer_in": 255960926,
    "transfer_out": 311949952,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-4",
    "size": 366932245162,
    "storage_cost": 366.932245162,
    "bw_cost": 945.47,
    "total_cost": 1312.402245162,
    "transfer_in": 162543387,
    "transfer_out": 27261455,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-5",
    "size": 366720789126,
    "storage_cost": 366.720789126,
    "bw_cost": 880.08,
    "total_cost": 1246.800789126,
    "transfer_in": 228731291,
    "transfer_out": 116688694,
    "chargify": "synced"
  },
  {
    "_id": "bucket-6",
    "size": 184096299903,
    "storage_cost": 184.09629990300002,
    "bw_cost": 428.53,
    "total_cost": 612.626299903,
    "transfer_in": 337711916,
    "transfer_out": 143932895,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-7",
    "size": 120118032943,
    "storage_cost": 120.118032943,
    "bw_cost": 783.37,
    "total_cost": 903.488032943,
    "transfer_in": 53696798,
    "transfer_out": 107032128,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-8",
    "size": 92824000290,
    "storage_cost": 92.82400029,
    "bw_cost": 436.93,
    "total_cost": 529.75400029,
    "transfer_in": 105581546,
    "transfer_out": 57611939,
    "chargify": "synced"
  },
  {
    "_id": "bucket-9",
    "size": 116615463902,
    "storage_cost": 116.615463902,
    "bw_cost": 790.68,
    "total_cost": 907.2954639019999,
    "transfer_in": 365488105,
    "transfer_out": 399135452,
    "chargify": "synced"
  },
  {
    "_id": "bucket-10",
    "size": 171385867952,
    "storage_cost": 171.385867952,
    "bw_cost": 236.02,
    "total_cost": 407.405867952,
    "transfer_in": 178134309,
    "transfer_out": 130783668,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-11",
    "size": 55642647341,
    "storage_cost": 55.642647341,
    "bw_cost": 830.01,
    "total_cost": 885.652647341,
    "transfer_in": 32260921,
    "transfer_out": 31714200,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-12",
    "size": 81551171697,
    "storage_cost": 81.551171697,
    "bw_cost": 976.86,
    "total_cost": 1058.411171697,
    "transfer_in": 224613504,
    "transfer_out": 119709958,
    "chargify": "synced"
  },
  {
    "_id": "bucket-13",
    "size": 144409602121,
    "storage_cost": 144.409602121,
    "bw_cost": 939.49,
    "total_cost": 1083.899602121,
    "transfer_in": 286813425,
    "transfer_out": 320786164,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-14",
    "size": 344638779320,
    "storage_cost": 344.63877932,
    "bw_cost": 274.77,
    "total_cost": 619.40877932,
    "transfer_in": 46554528,
    "transfer_out": 326513379,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-15",
    "size": 294065069493,
    "storage_cost": 294.065069493,
    "bw_cost": 552.54,
    "total_cost": 846.605069493,
    "transfer_in": 125468019,
    "transfer_out": 366529372,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-16",
    "size": 180419506218,
    "storage_cost": 180.419506218,
    "bw_cost": 439.82,
    "total_cost": 620.239506218,
    "transfer_in": 345929955,
    "transfer_out": 120241519,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-17",
    "size": 254433445685,
    "storage_cost": 254.433445685,
    "bw_cost": 197.96,
    "total_cost": 452.393445685,
    "transfer_in": 82613749,
    "transfer_out": 247013050,
    "chargify": "synced"
  },
  {
    "_id": "bucket-18",
    "size": 379260721228,
    "storage_cost": 379.260721228,
    "bw_cost": 458.58,
    "total_cost": 837.840721228,
    "transfer_in": 309130516,
    "transfer_out": 188678636,
    "chargify": "synced"
  },
  {
    "_id": "bucket-19",
    "size": 90563537176,
    "storage_cost": 90.56353717600001,
    "bw_cost": 163.57,
    "total_cost": 254.133537176,
    "transfer_in": 72110959,
    "transfer_out": 266563744,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-20",
    "size": 217797250004,
    "storage_cost": 217.797250004,
    "bw_cost": 255.88,
    "total_cost": 473.67725000400003,
    "transfer_in": 61939786,
    "transfer_out": 267682416,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-21",
    "size": 243049158845,
    "storage_cost": 243.04915884500002,
    "bw_cost": 717.39,
    "total_cost": 960.4391588450001,
    "transfer_in": 202466931,
    "transfer_out": 291831746,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-22",
    "size": 264347506760,
    "storage_cost": 264.34750676000004,
    "bw_cost": 693.46,
    "total_cost": 957.80750676,
    "transfer_in": 183184873,
    "transfer_out": 137230803,
    "chargify": "in_sync"
  },
  {
    "_id": "bucket-23",
    "size": 324076192386,
    "storage_cost": 324.076192386,
    "bw_cost": 712.44,
    "total_cost": 1036.516192386,
    "transfer_in": 224637482,
    "transfer_out": 174271917,
    "chargify": "synced"
  },
  {
    "_id": "bucket-24",
    "size": 316656096154,
    "storage_cost": 316.656096154,
    "bw_cost": 859.88,
    "total_cost": 1176.536096154,
    "transfer_in": 192468701,
    "transfer_out": 381607022,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-25",
    "size": 331255034043,
    "storage_cost": 331.255034043,
    "bw_cost": 671.9,
    "total_cost": 1003.155034043,
    "transfer_in": 256678545,
    "transfer_out": 316156903,
    "chargify": "synced"
  },
  {
    "_id": "bucket-26",
    "size": 257881298447,
    "storage_cost": 257.881298447,
    "bw_cost": 104.84,
    "total_cost": 362.721298447,
    "transfer_in": 47154953,
    "transfer_out": 310614794,
    "chargify": "not_synced"
  },
  {
    "_id": "bucket-27",
    "size": 151212762277,
    "storage_cost": 151.212762277,
    "bw_cost": 616.5,
    "total_cost": 767.712762277,
    "transfer_in": 357408451,
    "transfer_out": 281409199,
    "chargify": "not_synced"
  }
];

export const getS3Data = () => {
  return { success: true, data: mockData };
};
