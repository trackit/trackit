import React, { Component } from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';
import { Link } from 'react-router-dom';
import { formatPrice } from '../../common/formatters';


class TopUnusedComponent extends Component {
    getTotalItemsNB(unused) {
        let res = 0;

        if (unused.ec2)
            res += unused.ec2.values ? unused.ec2.values.length : 0;
        return res;
    }

    getPotentialsSavings(unused ) {
        let res = 0;

        if (unused.ec2 && unused.ec2.status && unused.ec2.values) {
            for (let i = 0; i < unused.ec2.values.length; i++) {
                const element = unused.ec2.values[i];
                res += element.cost;
            }
        }
        return res;
    }

    render() {
        let unused = [];
        const propsValues = this.props.unused;
        const ec2 = propsValues.ec2;

        if (ec2 && ec2.status && ec2.values) {
            for (let i = 0; i < ec2.values.length; i++) {
                const element = ec2.values[i];
                unused.push(
                    <tr key={element.id}>
                        <td className="badge-cell">
                            <span className="badge blue-bg">
                                EC2
                            </span>
                        </td>
                        <td>{element.tags.Name ? element.tags.Name : element.id}</td>
                        <td><strong>{formatPrice(element.cost)}</strong></td>
                        <td><i className="fa fa-microchip blue-color"></i> CPU avg is low : {element.cpuAverage.toFixed(1)}%</td>
                    </tr>
                );
            }
        }

        let noUnusedMessage;
        let table;
        if (!this.getTotalItemsNB(this.props.unused)) {
            noUnusedMessage= <h4 className="hl-panel-title no-resource-message">
                <i className="fa fa-check"></i>
                &nbsp;
                All good. No detected unused resources.
            </h4>;
        } else {
            table = (
                <div className="hl-table-wrapper">
                <table className="hl-table">
                    <colgroup>
                        <col style={{width: '10%' }}></col>
                        <col style={{width: '40%' }}></col>
                        <col style={{width: '15%' }}></col>
                        <col style={{width: '35%' }}></col>
                    </colgroup>  

                    <thead>
                        <tr>
                            <th></th>
                            <th>Name</th>
                            <th>Cost</th>
                            <th>Reason</th>
                        </tr>
                    </thead>
                    <tbody>
                        {unused}
                    </tbody>
                </table>
            </div>
            );
        }
      
        return (
            <div className="col-md-6">
                <div className="white-box hl-panel">
                    <h4 className="m-t-0 hl-panel-title">
                        {moment(this.props.date).format('MMM Y')} Unused resources
                        <br />
                        <span>({formatPrice(this.getPotentialsSavings(this.props.unused))}/mth potential savings)</span>
                    </h4>
                    <Link to="/app/resources" className="hl-details-link">
                        More details
                    </Link>
                    <hr className="m-b-0"/>
                    {noUnusedMessage}
                    {table}
                    <div className="clearfix"></div>
                </div>
            </div>
        );
    }
}

TopUnusedComponent.propTypes = {
    unused: PropTypes.object.isRequired,
    date: PropTypes.object.isRequired,
}

export default TopUnusedComponent;