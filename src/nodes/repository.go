package nodes

import (
	"database/sql"
	"log"
	"time"

	db "github.com/oci-hpc/slurm-cluster-api/src/database"
)

func upsertNodeStatus(node NodeInfo) {
	res := queryNodeStatus(node.Name)
	if res.Id == 0 {
		insertNodeStatus(node)
	} else {
		node.Id = res.Id
		updateNodeStatus(node)
	}
}

func updateNodeStatus(node NodeInfo) {
	node.LastSeenTime = time.Now()
	sqlString := `
		UPDATE t_node
		SET
			m_node_addr = :m_node_addr,
			m_node_hostname = :m_node_hostname,
			m_node_state = :m_node_state,
			m_port = :m_port,
			m_reason = :m_reason,
			m_reason_time = :m_reason_time,
			m_version = :m_version,
			m_last_seen_time = :m_last_seen_time
		WHERE id = :id;
	`
	dbConn := db.GetDbConnection()
	defer dbConn.Close()
	_, err := dbConn.Exec(
		sqlString,
		sql.Named("m_node_addr", node.NodeAddr),
		sql.Named("m_node_hostname", node.NodeHostname),
		sql.Named("m_node_state", node.NodeState),
		sql.Named("m_port", node.Port),
		sql.Named("m_reason", node.Reason),
		sql.Named("m_reason_time", node.ReasonTime),
		sql.Named("m_version", node.Version),
		sql.Named("m_last_seen_time", node.LastSeenTime),
		sql.Named("id", node.Id),
	)
	if err != nil {
		log.Printf("WARN: updateNodeStatus: " + err.Error())
	}
}

func insertNodeStatus(node NodeInfo) {
	node.LastSeenTime = time.Now()
	sqlString := `
		INSERT INTO t_node (
			m_cluster_name,
			m_cores,
			m_cpus,
			m_gres,
			m_name,
			m_node_addr,
			m_node_hostname,
			m_node_state,
			m_port,
			m_reason,
			m_reason_time,
			m_sockets,
			m_threads,
			m_version,
			m_last_seen_time 
		) values (
			:m_cluster_name,
			:m_cores,
			:m_cpus,
			:m_gres,
			:m_name,
			:m_node_addr,
			:m_node_hostname,
			:m_node_state,
			:m_port,
			:m_reason,
			:m_reason_time,
			:m_sockets,
			:m_threads,
			:m_version,
			:m_last_seen_time
		)
	`
	dbConn := db.GetDbConnection()
	defer dbConn.Close()
	_, err := dbConn.Exec(
		sqlString,
		sql.Named("m_cluster_name", node.ClusterName),
		sql.Named("m_cores", node.Cores),
		sql.Named("m_cpus", node.Cpus),
		sql.Named("m_gres", node.Gres),
		sql.Named("m_name", node.Name),
		sql.Named("m_node_addr", node.NodeAddr),
		sql.Named("m_node_hostname", node.NodeHostname),
		sql.Named("m_node_state", node.NodeState),
		sql.Named("m_port", node.Port),
		sql.Named("m_reason", node.Reason),
		sql.Named("m_reason_time", node.ReasonTime),
		sql.Named("m_sockets", node.Sockets),
		sql.Named("m_threads", node.Threads),
		sql.Named("m_version", node.Version),
		sql.Named("m_last_seen_time", node.LastSeenTime),
	)
	if err != nil {
		log.Printf("WARN: insertNodeStatus: " + err.Error())
	}
}

func queryNodeStatus(name string) (node NodeInfo) {
	sqlString := `
		SELECT 
			id,
			m_cluster_name,
			m_cores,
			m_cpus,
			m_gres,
			m_name,
			m_node_addr,
			m_node_hostname,
			m_node_state,
			m_port,
			m_reason,
			m_reason_time,
			m_sockets,
			m_threads,
			m_version,
			m_last_seen_time 
		FROM t_node
		WHERE m_name = :m_name;
	`
	dbConn := db.GetDbConnection()
	defer dbConn.Close()
	rows, err := dbConn.Query(sqlString, sql.Named("m_name", name))
	if err != nil {
		log.Printf("WARN: queryNodeStatus: " + err.Error())
	}
	var nodes []NodeInfo
	convertRowsToNodes(rows, &nodes)
	err = rows.Err()
	if err != nil {
		log.Printf("WARN: queryNodeStatus: " + err.Error())
	}
	if len(nodes) == 0 {
		return
	}
	return nodes[0]
}

func convertRowsToNodes(rows *sql.Rows, nodes *[]NodeInfo) {
	defer rows.Close()
	if rows == nil {
		return
	}
	for rows.Next() {
		var node NodeInfo
		err := rows.Scan(
			&node.Id,
			&node.ClusterName,
			&node.Cores,
			&node.Cpus,
			&node.Gres,
			&node.Name,
			&node.NodeAddr,
			&node.NodeHostname,
			&node.NodeState,
			&node.Port,
			&node.Reason,
			&node.ReasonTime,
			&node.Sockets,
			&node.Threads,
			&node.Version,
			&node.LastSeenTime,
		)
		if err != nil {
			log.Printf("WARN: convertRowsToNodes: " + err.Error())
		}
		*nodes = append(*nodes, node)
	}
}

func queryAllNodeStatus() (nodes []NodeInfo) {
	sqlString := `
		SELECT 
			id,
			m_cluster_name,
			m_cores,
			m_cpus,
			m_gres,
			m_name,
			m_node_addr,
			m_node_hostname,
			m_node_state,
			m_port,
			m_reason,
			m_reason_time,
			m_sockets,
			m_threads,
			m_version,
			m_last_seen_time 
		FROM t_node;
	`
	dbConn := db.GetDbConnection()
	defer dbConn.Close()
	rows, err := dbConn.Query(sqlString)
	if err != nil {
		log.Printf("WARN: queryNodeStatus: " + err.Error())
	}
	convertRowsToNodes(rows, &nodes)
	err = rows.Err()
	if err != nil {
		log.Printf("WARN: queryNodeStatus: " + err.Error())
	}
	return nodes
}
