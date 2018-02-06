import React from 'react';

// Return a value not negative or zero
export const noNeg = (value) => (value < 0 ? 0 : value);

export const capitalizeFirstLetter = (value) => (value.charAt(0).toUpperCase() + value.slice(1));

// Take bytes value and return formatted string value. Second param is optional floating number
export const formatBytes = (a,d = 2) => {if(0===a)return"0 Bytes";var c=1024,e=["Bytes","KB","MB","GB","TB","PB","EB","ZB","YB"],f=Math.floor(Math.log(a)/Math.log(c));return parseFloat((a/Math.pow(c,f)).toFixed(d))+" "+e[f]};
export const formatGigaBytes = (a,d = 2) => (formatBytes(a * Math.pow(1024,3), d));

export const formatPrice = (value, decimals = 2) => (<span><span className="dollar-sign">$</span>{parseFloat(value.toFixed(decimals)).toLocaleString()}</span>);

export const transformProducts = (data, filter, interval) => {
  if (filter === "all" && data.hasOwnProperty(interval))
    return [{
      key: "Total",
      values: Object.keys(data[interval]).map((date) => ([date, data[interval][date]]))
    }];
  else if (!data.hasOwnProperty(filter))
    return [];
  let dates = [];
  try {
    Object.keys(data[filter]).forEach((key) => {
      Object.keys(data[filter][key][interval]).forEach((date) => {
        if (dates.indexOf(date) === -1)
          dates.push(date);
      })
    });
    return Object.keys(data[filter]).map((key) => ({
      key: (key.length ? key : `No ${filter}`),
      values: dates.map((date) => ([date, data[filter][key][interval][date] || 0]))
    }));
  } catch (e) {
    return [];
  }
};

export const transformBuckets = (data) => {
  return  Object.keys(data).map((bucket) => ({
    key: bucket,
    values: [
      ["Bandwidth", data[bucket].BandwidthCost],
      ["Storage", data[bucket].StorageCost]
    ]
  }));
/*  return [{
    key: "Bandwidth",
    values: Object.keys(data).map((bucket) => ([bucket, data[bucket].BandwidthCost])).slice(0, 2)
  },{
    key: "Storage",
    values: Object.keys(data).map((bucket) => ([bucket, data[bucket].StorageCost])).slice(0, 2)
  }];*/
};
