import React from 'react';

const isUniqueInArray = (value, index, array) => (array.every((elem, idx) => (elem !== value || index === idx)));
const isString = (value) => (typeof value === "string");
const isNumber = (value) => (!isNaN(value));
const isPositiveInteger = (value) => (isNumber(value) && Number.isInteger(parseFloat(value)) && parseInt(value, 10) >= 0);
const isPositiveIntegerInRange = (min, max) => (value) => (isPositiveInteger(value) && value >= min && value <= max);
const isArray = (value) => (Array.isArray(value) && value.every(isUniqueInArray));
const isWithPositiveIntegerInRangeArray = (min, max) => (value) => (isArray(value) && value.every(isPositiveIntegerInRange(min, max)));
const isStringArray = (value) => (isArray(value) && value.every(isString));

const weekDays = {
  0: "Monday",
  1: "Tuesday",
  2: "Wednesday",
  3: "Thursday",
  4: "Friday",
  5: "Saturday",
  6: "Sunday",
};
const levels = {
  0: "Low",
  1: "Medium",
  2: "High",
  3: "Critical"
};

const dateFilters = {
  week_day: {
    pretty: "Days of the week",
    input: {
      format: "checkbox",
      values: weekDays
    },
    name: "Specific days of the week",
    icon: <i className="fa fa-calendar"/>,
    validation: isWithPositiveIntegerInRangeArray(0, 6),
    format: (value) =>(<span className="event-filter-value-weekday">{value.map((day, index) => (<strong key={index}>{weekDays[day]}</strong>))}</span>)
  },
  month_day: {
    pretty: "Days of the month",
    input: {
      format: "checkbox",
      values: Array(31).fill().map((_, index) => index + 1)
    },
    name: "Specific days of the month",
    icon: <i className="fa fa-calendar"/>,
    validation: isWithPositiveIntegerInRangeArray(0, 30),
    format: (value) =>(<span className="event-filter-value-monthday">{value.map((day, index) => (<strong key={index}>{day}</strong>))}</span>)
  }
};

const costFilters = {
  cost_min: {
    pretty: "Minimum cost",
    input: {
      format: "input",
      type: "number"
    },
    name: "Hide events with cost lower than a given value",
    icon: <i className="fa fa-usd"/>,
    validation: isNumber,
    format: (value) => (<span className="event-filter-value-number"><strong>{value}</strong></span>)
  },
  cost_max: {
    pretty: "Maximum cost",
    input: {
      format: "input",
      type: "number"
    },
    name: "Hide events with cost higher than a given value",
    icon: <i className="fa fa-usd"/>,
    validation: isNumber,
    format: (value) => (<span className="event-filter-value-number"><strong>{value}</strong></span>)
  },
  expected_cost_min: {
    pretty: "Minimum expected cost",
    input: {
      format: "input",
      type: "number"
    },
    name: "Hide events with expected cost lower than a given value",
    icon: <i className="fa fa-usd"/>,
    validation: isNumber,
    format: (value) => (<span className="event-filter-value-number"><strong>{value}</strong></span>)
  },
  expected_cost_max: {
    pretty: "Maximum expected cost",
    input: {
      format: "input",
      type: "number"
    },
    name: "Hide events with expected cost higher than a given value",
    icon: <i className="fa fa-usd"/>,
    validation: isNumber,
    format: (value) => (<span className="event-filter-value-number"><strong>{value}</strong></span>)
  }
};

export const filters = {
  ...dateFilters,
  ...costFilters,
  product: {
    pretty: "Product",
    input: {
      format: "array",
      type: "string"
    },
    name: "Filter by product name",
    validation: isStringArray,
    format: (value) =>(<span className="event-filter-value-products">{value.map((product, index) => (<strong key={index}>{product}</strong>))}</span>)
  },
  level: {
    pretty: "Level",
    input: {
      format: "checkbox",
      values: levels
    },
    name: "Filter by event level",
    validation: isWithPositiveIntegerInRangeArray(0, 3),
    format: (value) =>(<span className="event-filter-value-products">{value.map((level, index) => (<strong key={index}>{levels[level]}</strong>))}</span>)
  }
};

export const getFilterName = (rule) => ((rule !== null && filters.hasOwnProperty(rule) && filters[rule].hasOwnProperty("name")) ? filters[rule].name : "Unknown filter");
export const getFilterIcon = (rule) => ((rule !== null && filters.hasOwnProperty(rule) && filters[rule].hasOwnProperty("icon")) ? filters[rule].icon : <i className="fa fa-filter"/>);
export const checkFilterValue = (rule, value) => ((rule !== null && filters.hasOwnProperty(rule) && filters[rule].hasOwnProperty("validation")) ? filters[rule].validation(value) : true);
export const getFilterValue = (rule, value) => {
  if (!checkFilterValue(rule, value))
    return (<div className="event-filter-value-invalid"><strong>Invalid format</strong></div>);
  if (rule !== null && filters.hasOwnProperty(rule) && filters[rule].hasOwnProperty("format"))
    return filters[rule].format(value);
  return (<span className="event-filter-value-general">{value}</span>);
};
export const getFilterInput = (rule) => ((rule !== null && filters.hasOwnProperty(rule) && filters[rule].hasOwnProperty("input")) ? filters[rule].input : null);

export const showFilter = (filter) => {
  return (
    <div className="event-filter event-filter-summary">
      <div className="event-filter-summary-name">
        {getFilterIcon(filter.rule)}
        &nbsp;
        {getFilterName(filter.rule)}
      </div>
      <div className="event-filter-summary-value">
        Value : {getFilterValue(filter.rule, filter.data)}
      </div>
    </div>
  )
};
