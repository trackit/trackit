import React, {Component} from 'react';
import {connect} from 'react-redux';

import Components from '../../components';

import Spinner from 'react-spinkit';
import PropTypes from "prop-types";
import Actions from "../../actions";

import {formatPrice} from '../../common/formatters';

const Map = Components.AWS.Map.Map;
const TimerangeSelector = Components.Misc.TimerangeSelector;

const regions = {
  "sa-east-1": "Sao Paulo",
  "ca-central-1": "Canada",
  "cn-north-1": "Bejing",
  "cn-northwest-1": "Ningxia",
  "eu-central-1": "Frankfurt",
  "eu-west-1": "Ireland",
  "eu-west-2": "London",
  "eu-west-3": "Paris",
  "ap-northeast-1": "Tokyo",
  "ap-northeast-2": "Seoul",
  "ap-south-1": "Mumbai",
  "ap-southeast-1": "Singapore",
  "ap-southeast-2": "Sydney",
  "us-east-1": "North Virginia",
  "us-east-2": "Ohio",
  "us-west-1": "North California",
  "us-west-2": "Oregon",
  "global": "No specific region",
  "taxes": "Taxes"
};

const formatData = (costs) => {
  const data = {};
  let maxTotal = 0;
  Object.keys(regions).forEach((region) => {
    data[region] = {
      name: regions[region],
      total: 0,
      zones: {},
      opacity: 0
    };
    if (costs.hasOwnProperty("region"))
      Object.keys(costs.region).forEach((zone) => {
        if ((zone.startsWith(region) && region.length) || zone === region) {
          let total = 0;
          Object.keys(costs.region[zone].product).forEach((product) => {
            total += costs.region[zone].product[product];
          });
          data[region].zones[zone] = {
            total,
            products: costs.region[zone].product
          };
          data[region].total += total;
        }
      });
    if (data[region].total > maxTotal)
      maxTotal = data[region].total;
  });
  Object.keys(data).forEach((region) => {
    if (data[region].total) {
      const ratio = data[region].total / maxTotal;
      data[region].opacity = (ratio < 0.25 ? 0.5 : (ratio < 0.5 ? 0.7 : (ratio < 0.75 ? 0.9 : 1)));
    }
  });
  return data;
};

const regionDetails = (key, region, data, double, close) => {
  const zones = Object.keys(data.zones).map((zone, index) => (
    <div className="zone-item" key={index}>
      <div className="zone-name">
        {zone}
      </div>
      <div className="zone-products">
        {Object.keys(data.zones[zone].products).map((product, index) => (
          <div className="product-item" key={index}>
            {product} : {formatPrice(data.zones[zone].products[product])}
          </div>
        ))}
      </div>
    </div>
  ));
  const colWidth = (double ? "col-md-6" : "col-md-12");
  return (
    <div key={key} className={"region-details white-box " + colWidth}>
      <div className="header">
        <div className="close" onClick={close}>
          <i className="fa fa-times"/>
        </div>
      </div>
      <div className="region-name">
        <h3>{(region === "global" ? "Global Products" : (region === "taxes" ? "Taxes" : region))}</h3>
        <h4>{regions[region]}</h4>
      </div>
      <div className="region-info">
        <div>
          <div className="col-md-4 col-sm-4 p-t-15 p-b-15 br-sm br-md bb-xs info">
            <ul className="in-col">
              <li>
                <i className="fa fa-credit-card fa-2x blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {formatPrice(data.total)}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              total cost
            </h4>
          </div>
          <div className="col-md-4 col-sm-4 p-t-15 p-b-15 br-md bb-xs info">
            <ul className="in-col">
              <li>
                <i className="fa fa-th-list fa-2x blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {Object.keys(data.zones).map((zone) => (Object.keys(data.zones[zone].products).length)).reduce((a, b) => (a + b), 0)}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              products
            </h4>
          </div>
          <div className="col-md-4 col-sm-4 p-t-15 p-b-15 info">
            <ul className="in-col">
              <li>
                <i className="fa fa-globe fa-2x blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {Object.keys(data.zones).length}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              zones
            </h4>
          </div>
          <span className="clearfix"></span>
        </div>
        <div className="region-info-zones">
          {zones}
        </div>
      </div>
    </div>
  );
};

// ResourcesMapContainer Component
export class ResourcesMapContainer extends Component {

  constructor(props){
    super(props);
    this.state = {
      selected: [],
      data: {}
    };
    this.selectRegion = this.selectRegion.bind(this);
  }

  componentWillMount() {
    this.props.getCosts(this.props.dates.startDate, this.props.dates.endDate);
  }

  componentWillUnmount() {
    this.props.clearCosts();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts)
      nextProps.getCosts(nextProps.dates.startDate, nextProps.dates.endDate);
    else {
      this.setState({selected: []});
      if (nextProps.costs.status && nextProps.costs.hasOwnProperty("values"))
        this.setState({data: formatData(nextProps.costs.values)});
    }
  }

  selectRegion = (newSelected) => {
    const selected = this.state.selected;
    if (selected.indexOf(newSelected) !== -1)
      selected.splice(selected.indexOf(newSelected), 1);
    else {
      selected.push(newSelected);
      while (selected.length > 2)
        selected.shift() 
    }
    this.setState({selected});
  };

  unselectRegion = (region) => {
    const selected = this.state.selected;
    selected.splice(selected.indexOf(region), 1);
    this.setState({selected});
  };

  render() {
    const loading = (!this.props.costs || !this.props.costs.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

    const error = (this.props.costs && this.props.costs.status && this.props.costs.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">Data not available ({this.props.costs.error.message})</div>
    ) : null);

    const map = (!loading && !error ? (
      <Map data={this.state.data} selectRegion={this.selectRegion}/>
    ) : null);

    const emptySelection = (
      <div className="white-box">
        <div className="map-empty-selection">
          <i className="fa fa-map-o"/>
          &nbsp;
          Select a region to see more details (You can select up to 2 regions)
        </div>
      </div>
    );

    const selection = this.state.selected.map((item, index) => regionDetails(index, item, this.state.data[item], (this.state.selected.length === 2), this.unselectRegion.bind(this, item)));
    const selectionDetails = (!this.state.selected.length ? emptySelection : (
      <div className="row row-eq-height row-regions-details">
        {selection}
      </div>
    ));

    let badges;

    if (this.props.costs && this.props.costs.status) {
      badges = (
        <Components.AWS.Accounts.StatusBadges
          values={
            this.props.costs ? (
              this.props.costs.status ? this.props.costs.values : {}
            ) : {}
          }
        />
      );
    }
  

    return (
      <div className="container-fluid">

        <div className="clearfix white-box">
          <div className="inline-block">
            <h3 className="white-box-title no-padding inline-block">
              <i className="fa fa-globe"></i>
              &nbsp;
              Resources Map
              {badges}
            </h3>
          </div>
          <div className="inline-block pull-right">
            <TimerangeSelector
              startDate={this.props.dates.startDate}
              endDate={this.props.dates.endDate}
              setDatesFunc={this.props.setDates}
            />
          </div>
        </div>

        <div className="white-box">
          {loading || error || map}
        </div>

        {selectionDetails}

      </div>
    );
  }

}

ResourcesMapContainer.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  costs: PropTypes.object,
  dates: PropTypes.object,
  getCosts: PropTypes.func.isRequired,
  clearCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
  resetDates: PropTypes.func.isRequired,
  clearDates: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  costs: aws.map.values,
  dates: aws.map.dates,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getCosts: (begin, end) => {
    dispatch(Actions.AWS.Map.getCosts(begin, end));
  },
  clearCosts: () => {
    dispatch(Actions.AWS.Map.clearCosts());
  },
  setDates: (startDate, endDate) => {
    dispatch(Actions.AWS.Map.setDates(startDate, endDate));
  },
  resetDates: () => {
    dispatch(Actions.AWS.Map.resetDates());
  },
  clearDates: () => {
    dispatch(Actions.AWS.Map.clearDates());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesMapContainer);
