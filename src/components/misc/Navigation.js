import React, {Component} from 'react';
import { connect } from 'react-redux';
import { NavLink, withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import Actions from '../../actions';
import NavbarHeader from './NavbarHeader';

// Styling
import '../../styles/Navigation.css';

// Assets

export class Navigation extends Component {

  constructor() {
    super();
    this.state = {
      userMenuExpanded: false,
    };
  }

  componentDidMount() {
    this.props.getData(this.props.eventsDates.startDate, this.props.eventsDates.endDate);
  }

  getAbnormals(data) {
    const res = [];
    for (let i = 0; i < data.length; i++) {
      const element = data[i];
      if (element.abnormal) {
        res.push(element);
      }
    }
    return res.length ? res : null;
  }


  getEventsNumber(data) {
    let res = 0;
    for (var key in data) {
      if (data.hasOwnProperty(key)) {
          const dataSet = data[key];
          const abnormals = this.getAbnormals(dataSet);
          if (abnormals) {
            res += abnormals.length;
          }
      }
    }
    return res;
  }

  render() {

    let eventsBadge;
    if (this.props.events && this.props.events.status && this.props.events.values) {
      eventsBadge = <span className="menu-badge">{this.getEventsNumber(this.props.events.values)}</span>;
    }

    return(
      <div className="navigation">

        <NavbarHeader />

        <div className="navbar-default sidebar animated fadeInLeft" role="navigation">
          <div className="sidebar-head">
            <h3>
              <span className="open-close">
                <i className="fa fa-bars hidden-xs"></i>
              </span>
              <span className="hide-menu">
                Navigation
              </span>
            </h3>
          </div>

          <ul className="nav flex-column" id="side-menu">
            <li className="nav-item">
              <NavLink className="nav-link" exact to='/app' activeClassName="active">
                <i className="menu-icon fa fa-home"/>
                <span className="hide-menu">Home</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/dashboard' activeClassName="active">
                <i className="menu-icon fa fa-tachometer"/>
                <span className="hide-menu">Dashboard</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/costbreakdown' activeClassName="active">
                <i className="menu-icon fa fa-area-chart"/>
                <span className="hide-menu">Cost Breakdown</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/events' activeClassName="active">
                <div className="menu-icon inline-block">
                <i className="menu-icon fa fa-exclamation-triangle"/>
                {eventsBadge}
                </div>
                <span className="hide-menu">Events</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/plugins' activeClassName="active">
                <i className="menu-icon fa fa-check-square-o"/>
                <span className="hide-menu">Recommendations</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/tags' activeClassName="active">
                <i className="menu-icon fa fa-tags"/>
                <span className="hide-menu">Tags</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/map' activeClassName="active">
                <i className="menu-icon fa fa-map-o"/>
                <span className="hide-menu">Resources Map</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/resources' activeClassName="active">
                <i className="menu-icon fa fa-list-alt"/>
                <span className="hide-menu">AWS Resources</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/s3' activeClassName="active">
                <i className="menu-icon fa fa-bar-chart"/>
                <span className="hide-menu">AWS S3 Analytics</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/reports' activeClassName="active">
                <i className="menu-icon fa fa-file-text"/>
                <span className="hide-menu">AWS Reports</span>
              </NavLink>
            </li>
            {/*
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/optimize' activeClassName="active">
                <i className="menu-icon fa fa-area-chart"/>
                <span className="hide-menu">Compute Optimizer</span>
              </NavLink>
            </li>
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/ressources' activeClassName="active">
                <i className="menu-icon fa fa-pie-chart"/>
                <span className="hide-menu">Ressources Monitoring</span>
              </NavLink>
            </li>
            */}
            <li className="nav-item">
              <NavLink className="nav-link" to='/app/setup' activeClassName="active">
                <i className="menu-icon fa fa-cog"/>
                <span className="hide-menu">Setup</span>
              </NavLink>
            </li>

          </ul>
        </div>
      </div>
    );
  }

}

Navigation.propTypes = {
  mail: PropTypes.string,
  eventsDates: PropTypes.object.isRequired,
  events: PropTypes.object.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  mail: state.auth.mail,
  events: state.events.values,
  eventsDates: state.events.dates,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (begin, end) => {
    dispatch(Actions.Events.getData(begin, end));
  },
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Navigation));
