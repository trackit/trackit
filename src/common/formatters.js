import React from 'react';

// Return a value not negative or zero
export const noNeg = (value) => (value < 0 ? 0 : value);

export const capitalizeFirstLetter = (value) => (value.charAt(0).toUpperCase() + value.slice(1));

// Take bytes value and return formatted string value. Second param is optional floating number
export const formatBytes = (a,d = 2) => {if(0===a)return"0 Bytes";var c=1024,e=["Bytes","KB","MB","GB","TB","PB","EB","ZB","YB"],f=Math.floor(Math.log(a)/Math.log(c));return parseFloat((a/Math.pow(c,f)).toFixed(d))+""+e[f]};
export const formatGigaBytes = (a,d = 2) => (formatBytes(a * Math.pow(1024,3), d));

export const formatPrice = (value, decimals = 2) => (<span><span className="dollar-sign">$</span>{parseFloat(value.toFixed(decimals)).toLocaleString()}</span>);

export const formatPercent = (value, decimals = 2, style=true) => {
  const formattedValue = parseFloat(Math.abs(value).toFixed(decimals)).toLocaleString();
  if (!style)
    return (<span className="percent ">{formattedValue}<span className="percent-sign">%</span></span>);
  const color = (value > 0 ? "red-color" : (value < 0 ? "green-color" : ""));
  const bold = (Math.abs(value) >= 100 ? "percent-bold " : "");
  const sign = (value > 0 ? "+" : (value < 0 ? "-" : ""));
  return (<span className={"percent " + bold + color}>{sign + formattedValue}<span className="percent-sign">%</span></span>);
};

export const formatDate = (moment, precision) => {
  switch (precision) {
    case "year":
      return moment.format('Y');
    case "month":
      return moment.format('MMM Y');
    case "day":
    default:
      return moment.format('MMM Do Y');
  }
};

export const costBreakdown = {
  transformProductsBarChart: (data, filter, interval) => {
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
      dates.sort();
      return Object.keys(data[filter]).map((key) => ({
        key: (key.length ? key : `No ${filter}`),
        values: dates.map((date) => ([date, data[filter][key][interval][date] || 0]))
      }));
    } catch (e) {
      return [];
    }
  },
  transformProductsPieChart: (data, filter) => {
    if (!data.hasOwnProperty(filter))
      return [];
    return Object.keys(data[filter]).map((id) => ({
      key: id,
      value: data[filter][id]
    }));
  },
  getTotalPieChart: (data) => {
    let total = 0;
    if (Array.isArray(data))
      data.forEach((item) => {
        total += item.value;
      });
    return total;
  },
  transformCostDifferentiator: (data) => {
    let dates = [];
    const values = Object.keys(data).map((id) => {
      const itemValues = {};
      let previous = null;
      data[id].forEach((item) => {
        if (dates.indexOf(item.Date) === -1)
          dates.push(item.Date);
        let variation = item.PercentVariation;
        if (previous !== null && Math.abs(previous).toFixed(2) < 0.01 && Math.abs(item.Cost).toFixed(2) < 0.01)
          variation = 0;
        previous = item.Cost;
        itemValues[item.Date] = {
          cost: item.Cost,
          variation
        };
      });
      return {
        key: id,
        ...itemValues
      };
    });
    let previous = null;
    const total = {
      key: "Total"
    };
    dates.forEach((date) => {
      let cost = 0;
      values.forEach((value) => {
        cost += (value.hasOwnProperty(date) ? value[date].cost : 0);
      });
      const variation = (previous != null && previous !== 0 ? (cost - previous) / previous * 100 : 0);
      previous = cost;
      total[date] = {cost, variation};
    });
    return {dates, values, total};
  }
};

export const s3Analytics = {
  transformBuckets: (data) => {
    return Object.keys(data).map((bucket) => ({
      key: bucket,
      values: [
        ["Bandwidth", data[bucket].BandwidthCost],
        ["Storage", data[bucket].StorageCost]
      ]
    }))
  },
  transformBandwidthPieChart: (data) => {
    return Object.keys(data).map((bucket) => ({
      key: bucket,
      value: data[bucket].BandwidthCost
    }));
  },
  transformStoragePieChart: (data) => {
    return Object.keys(data).map((bucket) => ({
      key: bucket,
      value: data[bucket].StorageCost
    }));
  },
  transformRequestsPieChart: (data) => {
    return Object.keys(data).map((bucket) => ({
      key: bucket,
      value: data[bucket].RequestsCost
    }));
  },
  getTotalPieChart: (data) => {
    let total = 0;
    if (Array.isArray(data))
    data.forEach((item) => {
      total += item.value;
    });
    return total;
  }
};
