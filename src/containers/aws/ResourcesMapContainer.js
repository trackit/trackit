import React, {Component} from 'react';
import {connect} from 'react-redux';

import Components from '../../components';

import Spinner from 'react-spinkit';
import s3square from '../../assets/s3-square.png';
import PropTypes from "prop-types";
import Actions from "../../actions";

import {formatPrice} from '../../common/formatters';

const Panel = Components.Misc.Panel;
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
  "us-west-2": "Oregon"
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
        if (zone.startsWith(region)) {
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
      data[region].opacity = (ratio < 0.25 ? 0.4 : (ratio < 0.5 ? 0.6 : (ratio < 0.75 ? 0.8 : 1)));
    }
  });
  return data;
};

const regionDetails = (region, data, close) => {
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
  return (
    <div className="region-details">
      <div className="header">
        <div className="close" onClick={close}>
          <i className="fa fa-times"/>
        </div>
      </div>
      <div className="region-name">
        {region}
      </div>
      <div className="region-info">
        <div>
          <div className="col-md-3 col-md-offset-2 col-sm-4 p-t-15 p-b-15 br-sm br-md bb-xs">
            <ul className="in-col">
              <li>
                <i className="fa fa-dollar fa-2x green-color"/>
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
          <div className="col-md-3 col-sm-4 p-t-15 p-b-15 br-md bb-xs">
            <ul className="in-col">
              <li>
                <i className="fa fa-th-list fa-2x red-color"/>
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
          <div className="col-md-3 col-sm-4 p-t-15 p-b-15">
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
      selected: null,
      data: {}
    };
    this.selectRegion = this.selectRegion.bind(this);
    this.unselectRegion = this.unselectRegion.bind(this);
  }

  componentWillMount() {
    this.props.getCosts(this.props.dates.startDate, this.props.dates.endDate);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts)
      nextProps.getCosts(nextProps.dates.startDate, nextProps.dates.endDate);
    else {
      this.setState({selected: null});
      if (nextProps.costs.status && nextProps.costs.hasOwnProperty("values"))
        this.setState({data: formatData(nextProps.costs.values)});
    }
  }

  selectRegion = (selected) => {
    this.setState({selected});
  };

  unselectRegion = () => {
    this.setState({selected: null});
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
      <div className="map-empty-selection">
        <i className="fa fa-globe"/>
        &nbsp;
        Select a region to see more details
      </div>
    );

    return (
      <Panel>

        <div className="clearfix">
          <div className="inline-block">
            <h3 className="white-box-title no-padding inline-block">
              <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
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
        </div>

        <div>
          {loading || error || map}
        </div>

        <div>
          {!this.state.selected ?
            emptySelection :
            regionDetails(this.state.selected, this.state.data[this.state.selected], this.unselectRegion)
          }
        </div>

      </Panel>
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
    dispatch(Actions.AWS.Map.setDates(startDate, endDate))
  },
  resetDates: () => {
    dispatch(Actions.AWS.Map.resetDates())
  },
  clearDates: () => {
    dispatch(Actions.AWS.Map.clearDates())
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesMapContainer);
