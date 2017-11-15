import React from 'react';

// Return a value not negative or zero
export const noNeg = (value) => (value < 0 ? 0 : value);

export const capitalizeFirstLetter = (value) => (value.charAt(0).toUpperCase() + value.slice(1));

// Take bytes value and return formatted string value. Second param is optional floating number
export const formatBytes = (a,b) => {if(0===a)return"0 Bytes";var c=1024,d=b||2,e=["Bytes","KB","MB","GB","TB","PB","EB","ZB","YB"],f=Math.floor(Math.log(a)/Math.log(c));return parseFloat((a/Math.pow(c,f)).toFixed(d))+" "+e[f]}

export const formatPrice = (value, decimals = 2) => {return(<span><span className="dollar-sign">$</span>{parseFloat(value.toFixed(decimals)).toLocaleString()}</span>)}
