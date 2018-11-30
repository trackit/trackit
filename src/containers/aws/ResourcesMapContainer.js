import React, {Component} from 'react';
import {connect} from 'react-redux';
import Components from '../../components';
import Spinner from 'react-spinkit';
import PropTypes from "prop-types";
import Actions from "../../actions";
import {formatPrice} from '../../common/formatters';

const Map = Components.AWS.Map.Map;
const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;

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
  "": "No specific availability zone",
  "taxes": "Taxes"
};

let filters = {
  region: "Region",
  availabilityzone: "Availability zone"
};

const formatData = (costs, filter) => {
  const data = {};
  let maxTotal = 0;
  Object.keys(regions).forEach((region) => {
    data[region] = {
      name: regions[region],
      total: 0,
      zones: {},
      opacity: 0
    };
    const key = (costs.hasOwnProperty("region") ? "region" : (costs.hasOwnProperty("availabilityzone") ? "availabilityzone" : null));
    if (key)
      Object.keys(costs[key]).forEach((zone) => {
        if ((zone.startsWith(region) && region.length) || zone === region) {
          let total = 0;
          Object.keys(costs[key][zone].product).forEach((product) => {
            total += costs[key][zone].product[product];
          });
          data[region].zones[zone] = {
            total,
            products: costs[key][zone].product
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
  if (filter === "region")
    delete data[""];
  else
    delete data["global"];
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
  let title;
  switch (region) {
    case "global":
      title = "Global Products";
      break;
    case "taxes":
      title = "Taxes";
      break;
    case "":
      title = "Other products";
      break;
    default:
      title = region;
  }
  return (
    <div key={key} className={"region-details " + colWidth}>
      <div className="white-box">
        <div className="header">
          <div className="close" onClick={close}>
            <i className="fa fa-times"/>
          </div>
        </div>
        <div className="region-name">
          <h3>{title}</h3>
          <h4>{regions[region]}</h4>
        </div>
        <div className="region-info">
          <div className="row">
            <div className="col-md-6 col-sm-6 p-t-15 p-b-15 br-sm br-md bb-xs info">
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
            <div className="col-md-6 col-sm-6 p-t-15 p-b-15 info">
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
            <span className="clearfix"/>
          </div>
          <div className="region-info-zones">
            {zones}
          </div>
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
    this.props.getCosts(this.props.dates.startDate, this.props.dates.endDate, this.props.filter);
  }

  componentWillUnmount() {
    this.props.clearCosts();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates ||
      this.props.accounts !== nextProps.accounts ||
      this.props.filter !== nextProps.filter)
      nextProps.getCosts(nextProps.dates.startDate, nextProps.dates.endDate, nextProps.filter);
    else {
      this.setState({selected: []});
      if (nextProps.costs.status && nextProps.costs.hasOwnProperty("values"))
        this.setState({data: formatData(nextProps.costs.values, nextProps.filter)});
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
      <div className="col-md-12">
        <div className="white-box">
          <div className="map-empty-selection">
            <i className="fa fa-map-o"/>
            &nbsp;
            Select a region to see more details (You can select up to 2 regions)
          </div>
        </div>
      </div>
    );

    const selection = this.state.selected.map((item, index) => regionDetails(index, item, this.state.data[item], (this.state.selected.length === 2), this.unselectRegion.bind(this, item)));
    const selectionDetails = (!this.state.selected.length ? emptySelection : (
      <div className="row row-eq-height row-regions-details">
        {selection}
      </div>
    ));

    return (
      <div className="container-fluid">
        <div className="row">
          <div className="col-md-12">
            <div className="clearfix white-box">
              <div className="inline-block">
                <h3 className="white-box-title no-padding inline-block">
                  <i className="fa fa-globe"/>
                  &nbsp;
                  Resources Map
                </h3>
              </div>
              <div className="inline-block pull-right">
                <TimerangeSelector
                  startDate={this.props.dates.startDate}
                  endDate={this.props.dates.endDate}
                  setDatesFunc={this.props.setDates}
                />
              </div>
              <div className="inline-block pull-right">
                <Selector
                  values={filters}
                  selected={this.props.filter}
                  selectValue={this.props.setFilter}
                />
              </div>
            </div>

          </div>
        </div>
        
        <div className="white-box">
          {loading || error || map}
        </div>
        <div className="row">
          {selectionDetails}
        </div>
      </div>
    );
  }

}

ResourcesMapContainer.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  costs: PropTypes.object,
  dates: PropTypes.object,
  filter: PropTypes.string,
  getCosts: PropTypes.func.isRequired,
  clearCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
  resetDates: PropTypes.func.isRequired,
  clearDates: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  clearFilter: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  costs: aws.map.values,
  dates: aws.map.dates,
  filter: aws.map.filter,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getCosts: (begin, end, filter) => {
    dispatch(Actions.AWS.Map.getCosts(begin, end, filter));
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
  setFilter: (filter) => {
    dispatch(Actions.AWS.Map.setFilter(filter))
  },
  clearFilter: () => {
    dispatch(Actions.AWS.Map.clearFilter());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesMapContainer);
